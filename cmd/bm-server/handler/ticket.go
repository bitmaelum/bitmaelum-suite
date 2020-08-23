package handler

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"net/http"
)

/**
Tickets works as follows:

  We POST to either /incoming or /account/<hash>/incoming. The last one will use GetRemoteTicket and since this is
  authenticated, we always return a valid ticket. This ticket must be used to send a message through the API.

  If we POST to /incoming, we use the GetLocalTicket function and works differently.

    - if we post with (fromAddr, toAddr, and SubscriptionID) tuple, AND this tuple is found in the subscription list
      on the server, we automatically get a valid ticket. This ticket can be used directly for sending a message.
    - if we don't post a tuple, but only fromAddr and toAddr, we get a non-validated ticket in return. This ticket
      must be posted back with proof-of-work before it is exchanged into a valid ticket.
    - if we don't post a tuple, but a ticket ID, we check if the proof is ok. If so, we return a valid ticket.

  Note that tickets are always bound to a specific fromAddr and toAddr and only available for a specific lifetime.
*/

type ticketIn struct {
	FromAddr       string `json:"from_addr"`
	ToAddr         string `json:"to_addr"`
	SubscriptionID string `json:"subscription_id"`
	TicketID       string `json:"ticket_id"`
	Proof          uint64 `json:"proof_of_work"`
}

type httpError struct {
	err        string
	StatusCode int
}

func (e *httpError) Error() string {
	return e.err
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
	t := ticket.NewValidated("", *fromAddr, *toAddr, ticketBody.SubscriptionID)
	ticketRepo := container.GetTicketRepo()
	err = ticketRepo.Store(t)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
		return
	}

	logrus.Tracef("Generated valid ticket: %s", t.ID)

	// Send out validated ticket
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(t))
	return
}

// GetLocalTicket will try and retrieve a (valid) ticket so we can upload messages. It is only allowed to have a
// local destination address as this is used for local only (server-to-server communication)
func GetLocalTicket(w http.ResponseWriter, req *http.Request) {
	// Get ticket body from request or create a new ticket
	tckt, err := getTicketFromRequest(w, req)
	if err != nil {
		ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Check if the recipient address is valid
	err = validateLocalAddress(tckt)
	if err != nil {
		ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Check if we have a subscription tuple, if so, create valid ticket and return
	if tckt.SubscriptionID != "" {
		tckt, err = handleSubscription(tckt)
		if err != nil {
			ErrorOut(w, err.(*httpError).StatusCode, err.Error())
			return
		}

		outputTicket(tckt, w)
		return
	}

	// Ticket ID is not set, so we need to return a new unvalidated ticket
	if tckt.ID == "" {
		tckt = ticket.NewUnvalidated("", tckt.From, tckt.To, tckt.SubscriptionID)
		ticketRepo := container.GetTicketRepo()
		err = ticketRepo.Store(tckt)
		if err != nil {
			ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
			return
		}
		logrus.Tracef("Generated invalidated ticket: %s", tckt.ID)

		outputTicket(tckt, w)
		return
	}

	// Find the ticket in the repo
	ticketRepo := container.GetTicketRepo()
	tckt, err = ticketRepo.Fetch(tckt.ID)
	if err != nil {
		ErrorOut(w, http.StatusPreconditionFailed, "ticket not found")
		return
	}
	logrus.Tracef("Found ticket in repository: %s", tckt.ID)

	// Validate ticket if not done already
	if !tckt.Valid {
		tckt.Valid = tckt.Pow.HasDoneWork() && tckt.Pow.IsValid()
		if tckt.Valid {
			ticketRepo := container.GetTicketRepo()
			err = ticketRepo.Store(tckt)
			if err != nil {
				ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
				return
			}
			logrus.Tracef("Ticket proof-of-work validated: %s", tckt.ID)
		}
	}

	outputTicket(tckt, w)
	return
}

func outputTicket(tckt *ticket.Ticket, w http.ResponseWriter) {
	// Send out validated or invalidated ticket and status
	w.Header().Set("Content-Type", "application/json")
	if tckt.Valid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusPreconditionFailed)
	}
	_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(tckt))
}

func handleSubscription(tckt *ticket.Ticket) (*ticket.Ticket, error) {
	subscriptionRepo := container.GetSubscriptionRepo()

	// Check if we have the subscription stored
	sub := subscription.New(tckt.From, tckt.To, tckt.SubscriptionID)
	if !subscriptionRepo.Has(&sub) {
		return nil, &httpError{
			err:        "invalid subscription",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Subscription is valid, create a new validated ticket
	validTckt := tckt
	validTckt.Valid = true

	// Store the new validated ticket back in the repo
	ticketRepo := container.GetTicketRepo()
	err := ticketRepo.Store(validTckt)
	if err != nil {
		return nil, &httpError{
			err:        "can't save ticket on the server",
			StatusCode: http.StatusInternalServerError,
		}
	}

	logrus.Tracef("Validated ticket by subscription: %s %s", validTckt.ID, validTckt.SubscriptionID)
	return validTckt, nil
}

// validateLocalAddress checks if the recipient in the ticket is a local address. Returns error if not
func validateLocalAddress(tckt *ticket.Ticket) error {
	ar := container.GetAccountRepo()
	if !ar.Exists(tckt.To) {
		return &httpError{
			err:        "recipient isn't found on this server, and we don't support proxying",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil
}

// getTicketFromRequest returns a ticket based on the request input
func getTicketFromRequest(w http.ResponseWriter, req *http.Request) (*ticket.Ticket, error) {
	// Fetch ticket from ticket body
	ticketBody := &ticketIn{}
	err := DecodeBody(w, req.Body, ticketBody)
	if err != nil {
		return nil, &httpError{
			err:        err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// Validate from / to address
	fromAddr, err := address.NewHashFromHash(ticketBody.FromAddr)
	if err != nil {
		return nil, &httpError{
			err:        "Incorrect from address specified",
			StatusCode: http.StatusBadRequest,
		}
	}
	toAddr, err := address.NewHashFromHash(ticketBody.ToAddr)
	if err != nil {
		return nil, &httpError{
			err:        "Incorrect recipient address specified",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Create ticket with all info
	tckt := ticket.NewValidated(ticketBody.TicketID, *fromAddr, *toAddr, ticketBody.SubscriptionID)
	tckt.Pow.Proof = ticketBody.Proof

	return tckt, nil
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
