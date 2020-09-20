package handler

import (
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/davecgh/go-spew/spew"
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

type requestInfoType struct {
	From           *address.HashAddress
	To             *address.HashAddress
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

// GetClientToServerTicket will try and retrieves a valid ticket so we can upload messages. It is allowed to have a
// remote destination address as this is used by clients to upload messages for remote servers (client-to-server
// communication). Tickets here are always validated, since we need client authentication in order to fetch them.
func GetClientToServerTicket(w http.ResponseWriter, req *http.Request) {
	requestInfo, err := newFromRequest(req)
	if err != nil {
		ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Create new ticket, no need to validation
	// @TODO: subscription ID is empty here, but we probably want to fetch this directly from the server, not from the
	//  ticket body request (clients do not know anything about subscription ids)
	t := ticket.NewValidated(*requestInfo.From, *requestInfo.To, requestInfo.SubscriptionID)
	ticketRepo := container.GetTicketRepo()
	err = ticketRepo.Store(t)
	if err != nil {
		logrus.Trace("cannot save ticket: ", err)
		ErrorOut(w, http.StatusInternalServerError, "can't save ticket on the server")
		return
	}

	// Send out our validated ticket
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ticket.NewSimpleTicket(t))
	return
}

// GetServerToServerTicket will try and retrieve a (valid) ticket so we can upload messages. It is only allowed to have a
// local destination address as this is used for local only (server-to-server communication)
func GetServerToServerTicket(w http.ResponseWriter, req *http.Request) {
	// Get ticket body from request or create a new ticket
	requestInfo, err := newFromRequest(req)
	if err != nil {
		ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// @TODO: We need to take care of From and FromAddr (and To and ToAddr).. basically we can add anything here, \
	//  without us checking if they are actually hashes. Pretty much everything is accepted, which can cause conflict
	//  on the server.
	// For now: a hash SHOULD be verified

	logrus.Tracef("INFO FROM REQUEST: %#v", requestInfo)


	// Check if the recipient address is valid
	err = validateLocalAddress(*requestInfo.To)
	if err != nil {
		ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Check if we have a subscription tuple, if so, create valid ticket and return
	if requestInfo.SubscriptionID != "" {
		tckt, err := handleSubscription(requestInfo)
		if err != nil {
			ErrorOut(w, err.(*httpError).StatusCode, err.Error())
			return
		}

		outputTicket(tckt, w)
		return
	}

	// Ticket ID is not set, so we need to return a new unvalidated ticket
	if requestInfo.TicketID == "" {
		logrus.Trace("empty tckt.ID, creating new unvalidated ticket...")

		// Create new unvalidated ticket
		tckt := ticket.NewUnvalidated(*requestInfo.From, *requestInfo.To, requestInfo.SubscriptionID)

		// Store ticket
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
	tckt, err := ticketRepo.Fetch(requestInfo.TicketID)
	if err != nil {
		ErrorOut(w, http.StatusPreconditionFailed, "ticket not found")
		return
	}
	logrus.Tracef("Found ticket in repository: %s", tckt.ID)

	// Validate ticket if not done already
	if !tckt.Valid {
		// Set proof and check if it's valid
		tckt.Proof.Proof = requestInfo.Proof
		tckt.Valid = tckt.Proof.HasDoneWork() && tckt.Proof.IsValid()
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

func handleSubscription(requestInfo *requestInfoType) (*ticket.Ticket, error) {
	subscriptionRepo := container.GetSubscriptionRepo()

	// Check if we have the subscription stored
	sub := subscription.New(*requestInfo.From, *requestInfo.To, requestInfo.SubscriptionID)
	if !subscriptionRepo.Has(&sub) {
		return nil, &httpError{
			err:        "invalid subscription",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Subscription is valid, create a new validated ticket
	tckt := ticket.NewValidated(*requestInfo.From, *requestInfo.To, requestInfo.SubscriptionID)

	// Store the new validated ticket back in the repo
	ticketRepo := container.GetTicketRepo()
	err := ticketRepo.Store(tckt)
	if err != nil {
		return nil, &httpError{
			err:        "can't save ticket on the server",
			StatusCode: http.StatusInternalServerError,
		}
	}

	logrus.Tracef("Validated ticket by subscription: %s %s", tckt.ID, tckt.SubscriptionID)
	return tckt, nil
}

// validateLocalAddress checks if the recipient in the ticket is a local address. Returns error if not
func validateLocalAddress(addr address.HashAddress) error {
	ar := container.GetAccountRepo()
	if ar.Exists(addr) {
		return nil
	}

	return &httpError{
		err:        "recipient isn't found on this server, and we don't support proxying",
		StatusCode: http.StatusBadRequest,
	}
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

// newFromRequest returns a request info object based on the request body
func newFromRequest(req *http.Request) (*requestInfoType, error) {
	// Fetch all info from request
	requestInfo := &requestInfoType{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestInfo)
	if err != nil {
		return nil, &httpError{
			err:        "Malformed JSON: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// Validate from / to address
	requestInfo.From, err = address.NewHashFromHash(requestInfo.FromAddr)
	if err != nil {
		logrus.Trace("cannot create address: ", err)

		return nil, &httpError{
			err:        "Incorrect from address specified",
			StatusCode: http.StatusBadRequest,
		}
	}

	requestInfo.To, err = address.NewHashFromHash(requestInfo.ToAddr)
	if err != nil {
		logrus.Trace("cannot create address: ", err)

		return nil, &httpError{
			err:        "Incorrect to address specified",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Return info
	return requestInfo, nil
}
