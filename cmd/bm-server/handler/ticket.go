package handler

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ticketIn struct {
	FromAddr       string `json:"from_addr"`
	ToAddr         string `json:"to_addr"`
	SubscriptionID string `json:"subscription_id"`
	TicketID       string `json:"ticket_id"`
	Proof          uint64 `json:"proof_of_work"`
}

// GetRemoteTicket will try and retrieve a (valid) ticket so we can upload messages. It is allowed to have a
// remote destination address as this is used by clients to upload messages for remote servers (client-to-server communication)
func GetRemoteTicket(w http.ResponseWriter, req *http.Request) {
	ticketBody := &ticketIn{}
	err := DecodeBody(w, req.Body, ticketBody)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate from / to address
	fromAddr, err := address.NewHashFromHash(ticketBody.FromAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect from address specified")
		return
	}
	toAddr, err := address.NewHashFromHash(ticketBody.ToAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect to address specified")
		return
	}

	// Create new ticket, no need to validation
	t := ticket.NewValid(*fromAddr, *toAddr, ticketBody.SubscriptionID)
	ticketRepo := container.GetTicketRepo()
	err = ticketRepo.Store(t)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
		return
	}

	logrus.Tracef("Generated ticket: %s %v", t.ID, t.Valid)

	// Send out validated ticket
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(t))
	return
}

// GetLocalTicket will try and retrieve a (valid) ticket so we can upload messages. It is only allowed to have a
// local destination address as this is used for local only (server-to-server communication)
func GetLocalTicket(w http.ResponseWriter, req *http.Request) {
	ticketBody := &ticketIn{}
	err := DecodeBody(w, req.Body, ticketBody)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	var t *ticket.Ticket
	ticketRepo := container.GetTicketRepo()

	if ticketBody.TicketID != "" {
		t, err = ticketRepo.Fetch(ticketBody.TicketID)
		if err != nil {
			ErrorOut(w, http.StatusPreconditionFailed, "ticket not found")
			return
		}
		logrus.Tracef("Found ticket in repository: %s", ticketBody.TicketID)

		// Set info from ticket into body. It might overwrite, but ticket is leading..
		ticketBody.FromAddr = t.From.String()
		ticketBody.ToAddr = t.To.String()
		ticketBody.SubscriptionID = t.SubscriptionID

		// Set proof, we will check later if the proof is actually correct
		if ticketBody.Proof > 0 {
			t.Pow.Proof = ticketBody.Proof
		}
	}

	// Validate from / to address
	fromAddr, err := address.NewHashFromHash(ticketBody.FromAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect from address specified")
		return
	}
	toAddr, err := address.NewHashFromHash(ticketBody.ToAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect to address specified")
		return
	}

	// Check if to address is a local address
	ar := container.GetAccountRepo()
	if !ar.Exists(*toAddr) {
		ErrorOut(w, http.StatusBadRequest, "recipient isn't found on this server, and we don't support proxying")
		return
	}

	// Check if we have a subscription tuple, if so, create valid ticket and return
	subscriptionRepo := container.GetSubscriptionRepo()
	sub := subscription.New(*fromAddr, *toAddr, ticketBody.SubscriptionID)
	if subscriptionRepo.Has(&sub) {
		t := ticket.NewValid(*fromAddr, *toAddr, ticketBody.SubscriptionID)
		err = ticketRepo.Store(t)
		if err != nil {
			ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
			return
		}

		logrus.Tracef("Generated ticket: %s %v", t.ID, t.Valid)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(t))
		return
	}

	// Get ticket provided, or create a new ticket if not found
	if t == nil {
		t = ticket.New(*fromAddr, *toAddr, ticketBody.SubscriptionID)
		err = ticketRepo.Store(t)
		if err != nil {
			ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
			return
		}
		logrus.Tracef("Generated ticket: %s (valid: %v)", t.ID, t.Valid)
	}

	// Check if proof of work is done, and save accordingly
	if !t.Valid {
		t.Valid = t.Pow.HasDoneWork() && t.Pow.IsValid()
		if t.Valid {
			err = ticketRepo.Store(t)
			if err != nil {
				ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
				return
			}
		}

		logrus.Tracef("Ticket proof-of-work validated: %s (valid: %v)", t.ID, t.Valid)
	}

	// Send out validated or invalidated ticket
	w.Header().Set("Content-Type", "application/json")
	if t.Valid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusPreconditionFailed)
	}
	_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(t))
	return
}

// fetchTicketHeader returns a valid ticket as found in the request header, or err when no valid header or ticket is found.
func fetchTicketHeader(req *http.Request) (*ticket.Ticket, error) {
	ticketID := req.Header.Get(ticket.TicketHeader)

	ticketRepo := container.GetTicketRepo()
	t, err := ticketRepo.Fetch(ticketID)
	if err != nil {
		return nil, err
	}

	// Only return valid tickets
	if t.Valid == false {
		return nil, errors.New("invalid ticket")
	}

	logrus.Tracef("Valid ticket found: %s", t.ID)
	return t, nil
}
