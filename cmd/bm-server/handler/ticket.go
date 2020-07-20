package handler

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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
	body, err := readTicketBody(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate from / to address
	fromAddr, err := address.NewHashFromHash(body.FromAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect from address specified")
		return
	}
	toAddr, err := address.NewHashFromHash(body.ToAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect to address specified")
		return
	}

	// Create new ticket, no need to validation
	t := ticket.NewValid(*fromAddr, *toAddr, body.SubscriptionID)
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
	body, err := readTicketBody(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate from / to address
	fromAddr, err := address.NewHashFromHash(body.FromAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect from address specified")
		return
	}
	toAddr, err := address.NewHashFromHash(body.ToAddr)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "Incorrect to address specified")
		return
	}

	// Check if to address is a local address
	rs := container.GetResolveService()
	res, err := rs.Resolve(*toAddr)
	if err != nil || !res.IsLocal() {
		ErrorOut(w, http.StatusBadRequest, "destination isn't local, and we don't support proxying")
	}

	ticketRepo := container.GetTicketRepo()

	// Check if we have a subscription tuple, if so, create valid ticket and return
	subscriptionRepo := container.GetSubscriptionRepo()
	sub := subscription.New(*fromAddr, *toAddr, body.SubscriptionID)
	if subscriptionRepo.Has(&sub) {
		t := ticket.NewValid(*fromAddr, *toAddr, body.SubscriptionID)
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
	var t *ticket.Ticket
	if body.TicketID != "" {
		t, err = ticketRepo.Fetch(body.TicketID)
		if err != nil {
			ErrorOut(w, http.StatusPreconditionFailed, "ticket not found")
			return
		}
		logrus.Tracef("Found ticket in repository: %s", body.TicketID)
	}
	if t == nil {
		t = ticket.New(*fromAddr, *toAddr, body.SubscriptionID)
		err = ticketRepo.Store(t)
		if err != nil {
			ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
			return
		}
	}

	// Set proof, we will check later if the proof is actually correct
	if body.Proof > 0 {
		t.Pow.Proof = body.Proof
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

		logrus.Tracef("Ticket proof-of-work validated: %s %v", t.ID, t.Valid)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPreconditionFailed)
		_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(t))
		return
	}

	logrus.Tracef("Generated ticket: %s %v", t.ID, t.Valid)

	// Send out validated or invalidated ticket
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

// readTicketBody will read the incoming request body and fills a ticketIn struct
func readTicketBody(body io.ReadCloser) (*ticketIn, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	ticketBody := &ticketIn{}
	err = json.Unmarshal(data, &ticketBody)
	if err != nil {
		return nil, errors.New("incorrect body specified")
	}

	return ticketBody, nil
}
