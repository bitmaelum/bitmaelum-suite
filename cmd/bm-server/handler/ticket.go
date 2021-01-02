// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal/subscription"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/internal/work"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	errExpiredTicket  = errors.New("expired ticket")
	errInvalidTicket  = errors.New("invalid ticket")
	errCantSaveTicket = errors.New("can't save ticket on the server")
)

// RequestType is a request for a ticket
type RequestType struct {
	Sender         hash.Hash `json:"sender"`
	Recipient      hash.Hash `json:"recipient"`
	SubscriptionID string    `json:"subscription_id"`
	Preference     []string  `json:"preference"`
}

type httpError struct {
	err        string
	StatusCode int
}

// GetClientToServerTicket will try and retrieves a valid ticket so we can upload messages. It is allowed to have a
// remote destination address as this is used by clients to upload messages for remote servers (client-to-server
// communication). Tickets here are always validated, since we need client authentication in order to fetch them.
func GetClientToServerTicket(w http.ResponseWriter, req *http.Request) {
	requestInfo, err := newFromRequest(req)
	if err != nil {
		httputils.ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Create new ticket, no need to validation
	t := ticket.New(requestInfo.Sender, requestInfo.Recipient, requestInfo.SubscriptionID)
	t.Valid = true

	// Add authentication key if present in the request
	if IsAuthKeyAuthenticated(req) {
		logrus.Trace("adding authentication key to ticket")
		t.AuthKey = GetAuthKey(req).Fingerprint
	}

	ticketRepo := container.Instance.GetTicketRepo()
	err = ticketRepo.Store(t)
	if err != nil {
		logrus.Trace("cannot save ticket: ", err)
		httputils.ErrorOut(w, http.StatusInternalServerError, errCantSaveTicket.Error())
		return
	}

	// Send out our validated ticket
	_ = httputils.JSONOut(w, http.StatusOK, t)
}

// GetServerToServerTicket will try and retrieve a (valid) ticket so we can upload messages. It is only allowed to have a
// local destination address as this is used for local only (server-to-server communication)
func GetServerToServerTicket(w http.ResponseWriter, req *http.Request) {
	// Get ticket body from request or create a new ticket
	requestInfo, err := newFromRequest(req)
	if err != nil {
		httputils.ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Check if the recipient address is known locally (we don't support proxy)
	err = validateLocalAddress(requestInfo.Recipient)
	if err != nil {
		httputils.ErrorOut(w, err.(*httpError).StatusCode, err.Error())
		return
	}

	// Check if we have a subscription tuple, if so, create valid ticket and return
	if requestInfo.SubscriptionID != "" {
		t, err := handleSubscription(requestInfo)
		if err != nil {
			httputils.ErrorOut(w, err.(*httpError).StatusCode, err.Error())
			return
		}

		outputTicket(t, w)
		return
	}

	// Create new unvalidated ticket
	t := ticket.New(requestInfo.Sender, requestInfo.Recipient, requestInfo.SubscriptionID)

	// Add work based on the preference
	pw, err := work.GetPreferredWork(requestInfo.Preference)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "cannot create work for ticket")
		return
	}

	t.Work = &ticket.WorkType{
		Type: pw.GetName(),
		Data: pw,
	}

	// Store ticket
	ticketRepo := container.Instance.GetTicketRepo()
	err = ticketRepo.Store(t)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, errCantSaveTicket.Error())
		return
	}
	logrus.Tracef("Generated invalidated ticket: %s", t.ID)

	outputTicket(t, w)
}

// ValidateTicket will validate a ticket response
func ValidateTicket(w http.ResponseWriter, req *http.Request) {
	// Find the ticket in the repo
	ticketRepo := container.Instance.GetTicketRepo()
	t, err := ticketRepo.Fetch(mux.Vars(req)["ticket"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, "ticket not found")
		return
	}

	if t.Valid {
		// Ticket already valid. No need to check again
		outputTicket(t, w)
		return
	}

	// Check if ticket has work attached (it should)
	if t.Work == nil {
		httputils.ErrorOut(w, http.StatusExpectationFailed, "ticket does not contain work")
		return
	}

	// Read (json) body
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		httputils.ErrorOut(w, http.StatusExpectationFailed, "missing body")
		return
	}
	_ = req.Body.Close()

	// Try and validate the work
	t.Valid = t.Work.Data.ValidateWork(data)
	if t.Valid {
		t.Expiry = time.Now().Add(1800 * time.Second) // @TODO: Totally arbitrary
		t.Work = nil

		err = ticketRepo.Store(t)
		if err != nil {
			httputils.ErrorOut(w, http.StatusInternalServerError, errCantSaveTicket.Error())
			return
		}
		logrus.Tracef("Ticket proof-of-work validated: %s", t.ID)
	}

	outputTicket(t, w)
}

// output ticket info
func outputTicket(t *ticket.Ticket, w http.ResponseWriter) {
	_ = httputils.JSONOut(w, http.StatusOK, t)
}

func handleSubscription(requestInfo *RequestType) (*ticket.Ticket, error) {
	subscriptionRepo := container.Instance.GetSubscriptionRepo()

	// Check if we have the subscription stored
	sub := subscription.New(requestInfo.Sender, requestInfo.Recipient, requestInfo.SubscriptionID)
	if !subscriptionRepo.Has(&sub) {
		return nil, &httpError{
			err:        "invalid subscription",
			StatusCode: http.StatusBadRequest,
		}
	}

	// Subscription is valid, create a new validated ticket
	t := ticket.New(requestInfo.Sender, requestInfo.Recipient, requestInfo.SubscriptionID)
	t.Valid = true

	// Store the new validated ticket back in the repo
	ticketRepo := container.Instance.GetTicketRepo()
	err := ticketRepo.Store(t)
	if err != nil {
		return nil, &httpError{
			err:        errCantSaveTicket.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	logrus.Tracef("Validated ticket by subscription: %s %s", t.ID, t.SubscriptionID)
	return t, nil
}

// validateLocalAddress checks if the recipient in the ticket is a local address. Returns error if not
func validateLocalAddress(addr hash.Hash) error {
	ar := container.Instance.GetAccountRepo()
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

	ticketRepo := container.Instance.GetTicketRepo()
	t, err := ticketRepo.Fetch(ticketID)
	if err != nil {
		return nil, err
	}

	// Only return valid tickets
	if !t.Valid {
		return nil, errInvalidTicket
	}
	if t.Expired() {
		return nil, errExpiredTicket
	}

	logrus.Tracef("Valid ticket found: %s", t.ID)
	return t, nil
}

// newFromRequest returns a request info object based on the request body
func newFromRequest(req *http.Request) (*RequestType, error) {
	// Fetch all info from request
	requestInfo := &RequestType{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestInfo)
	if err != nil {
		return nil, &httpError{
			err:        "Malformed JSON: " + err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// Return info
	return requestInfo, nil
}

func (e *httpError) Error() string {
	return e.err
}
