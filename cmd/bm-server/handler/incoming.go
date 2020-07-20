package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

/*
 * Incoming will accept incoming messages from both clients and servers. It is based on (valid) tickets.
 */

// IncomingMessageHeader deals with uploading message headers
func IncomingMessageHeader(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		ErrorOut(w, http.StatusUnauthorized, "invalid ticket id or ticket not valid")
		return
	}

	// Read header from request body
	header, err := readHeaderFromBody(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "invalid header posted")
		return
	}

	// Verify from/to header with the ticket info
	if header.From.Addr.String() != t.From.String() || header.To.Addr.String() != t.To.String() {
		ErrorOut(w, http.StatusBadRequest, "header from/to address do not match the ticket")
		return
	}

	// Save request
	err = message.StoreHeader(t.ID, header)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message header")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("header saved"))
	return
}

// IncomingMessageCatalog deals with uploading message catalogs
func IncomingMessageCatalog(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		ErrorOut(w, http.StatusUnauthorized, "invalid ticket id or ticket not valid")
		return
	}

	err = message.StoreCatalog(t.ID, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message catalog")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved catalog"))
	return
}

// IncomingMessageBlock deals with uploading message blocks
func IncomingMessageBlock(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		ErrorOut(w, http.StatusUnauthorized, "invalid ticket id or ticket not valid")
		return
	}

	vars := mux.Vars(req)
	id := vars["id"]

	err = message.StoreBlock(t.ID, id, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message block")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved message block"))
	return
}

// CompleteIncoming is called whenever everything from a message has been uploaded and can be actually send
func CompleteIncoming(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		ErrorOut(w, http.StatusUnauthorized, "invalid ticket id or ticket not valid")
		return
	}

	// queue the message for processing
	processor.QueueIncomingMessage(t.ID)

	// Remove ticket
	ticketRepo := container.GetTicketRepo()
	ticketRepo.Remove(t.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(StatusOk("message accepted"))
	return
}

// DeleteIncoming is called whenever we want to completely remove a message by user request
func DeleteIncoming(w http.ResponseWriter, req *http.Request) {
	// Check ticket
	t, err := fetchTicketHeader(req)
	if err != nil {
		ErrorOut(w, http.StatusUnauthorized, "invalid ticket id or ticket not valid")
		return
	}

	if !message.IncomingPathExists(t.ID, "") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	err = message.RemoveMessage(message.SectionIncoming, t.ID)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while deleting outgoing message")
		return
	}

	// Remove ticket
	ticketRepo := container.GetTicketRepo()
	ticketRepo.Remove(t.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("message removed"))
	return
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
