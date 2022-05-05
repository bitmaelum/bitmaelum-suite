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
	"io"
	"io/ioutil"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
)

const (
	invalidTicketID string = "invalid ticket id or ticket not valid"
)

/*
 * Incoming will accept incoming messages from both clients and servers. It is based on (valid) tickets.
 */

// IncomingMessageHeader deals with uploading message headers
func IncomingMessageHeader(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		httputils.ErrorOut(w, http.StatusUnauthorized, invalidTicketID)
		return
	}

	// Read header from request body
	header, err := readHeaderFromBody(req.Body)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "invalid header posted")
		return
	}

	// Verify from/to header with the ticket info
	if header.From.Addr.String() != t.Sender.String() || header.To.Addr.String() != t.Recipient.String() {
		httputils.ErrorOut(w, http.StatusBadRequest, "header from/to address do not match the ticket")
		return
	}

	// If we are sending on behalf of another account, we need to add additional authorization information if we are originating
	if header.From.SignedBy == message.SignedByTypeAuthorized && IsOriginServer(header.From.Addr) {
		if t.AuthKey == "" {
			httputils.ErrorOut(w, http.StatusInternalServerError, "onbehalf signing, but token is not fetched with authentication")
			return
		}

		r := container.Instance.GetAuthKeyRepo()
		k, err := r.Fetch(t.AuthKey)
		if err != nil {
			httputils.ErrorOut(w, http.StatusInternalServerError, "error while fetching auth key")
			return
		}

		header.AuthorizedBy = &message.AuthorizedByType{
			PublicKey: k.PublicKey,
			Signature: k.Signature,
		}
	}

	// Add server signature if we are the originating server
	if IsOriginServer(header.From.Addr) {
		// Add a server signature to the header, so we know this server is the origin of the message
		err = message.SignServerHeader(header)
		if err != nil {
			httputils.ErrorOut(w, http.StatusInternalServerError, "error while signing incoming message")
			return
		}
	}

	// Save request
	err = message.StoreHeader(t.ID, header)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "error while storing message header")
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, httputils.StatusOk("saved message header"))
}

// IsOriginServer will return true when the sender of the header originates from this server. False otherwise
func IsOriginServer(addr hash.Hash) bool {
	rs := container.Instance.GetResolveService()
	info, err := rs.ResolveAddress(addr)
	if err != nil {
		return false
	}

	return info.RoutingID == config.Routing.RoutingID
}

// IncomingMessageCatalog deals with uploading message catalogs
func IncomingMessageCatalog(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		httputils.ErrorOut(w, http.StatusUnauthorized, invalidTicketID)
		return
	}

	err = message.StoreCatalog(t.ID, req.Body)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "error while storing message catalog")
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, httputils.StatusOk("saved message catalog"))
}

// IncomingMessageBlock deals with uploading message blocks
func IncomingMessageBlock(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		httputils.ErrorOut(w, http.StatusUnauthorized, invalidTicketID)
		return
	}

	vars := mux.Vars(req)
	messageID := vars["message"]

	err = message.StoreBlock(t.ID, messageID, req.Body)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "error while storing message block")
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, httputils.StatusOk("saved message block"))
}

// IncomingMessageAttachment deals with uploading message attachments
func IncomingMessageAttachment(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		httputils.ErrorOut(w, http.StatusUnauthorized, invalidTicketID)
		return
	}

	vars := mux.Vars(req)
	messageID := vars["message"]

	err = message.StoreAttachment(t.ID, messageID, req.Body)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "error while storing message attachment")
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, httputils.StatusOk("saved message attachment"))
}

// CompleteIncoming is called whenever everything from a message has been uploaded and can be actually send
func CompleteIncoming(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		httputils.ErrorOut(w, http.StatusUnauthorized, invalidTicketID)
		return
	}

	// queue the message for processing
	processor.QueueIncomingMessage(t.ID)

	// Remove ticket
	ticketRepo := container.Instance.GetTicketRepo()
	ticketRepo.Remove(t.ID)

	_ = httputils.JSONOut(w, http.StatusAccepted, httputils.StatusOk("message accepted"))
}

// DeleteIncoming is called whenever we want to completely remove a message by user request
func DeleteIncoming(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		httputils.ErrorOut(w, http.StatusUnauthorized, invalidTicketID)
		return
	}

	if !message.IncomingPathExists(t.ID, "") {
		httputils.ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	err = message.RemoveMessage(message.SectionIncoming, t.ID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusInternalServerError, "error while deleting outgoing message")
		return
	}

	// Remove ticket
	ticketRepo := container.Instance.GetTicketRepo()
	ticketRepo.Remove(t.ID)

	_ = httputils.JSONOut(w, http.StatusOK, httputils.StatusOk("message removed"))
}

func readHeaderFromBody(body io.ReadCloser) (*message.Header, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	header := &message.Header{}
	err = json.Unmarshal(data, &header)
	if err != nil {
		return nil, err
	}

	return header, nil
}
