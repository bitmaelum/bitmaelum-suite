package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

// Upload is when a CLIENT sends a message to the SERVER. This is an authenticated action so there is no need for
// any proof-of-work in this case. The client can only send to the server holding their account.

// UploadMessageHeader deals with uploading message headers
func UploadMessageHeader(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	// Did we already upload the header?
	if message.IncomingPathExists(msgID, "header.json") {
		ErrorOut(w, http.StatusConflict, "header already uploaded")
		return
	}

	// Read header from request body
	header, err := readHeaderFromBody(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "invalid header")
		return
	}

	// Save request
	err = message.StoreHeader(msgID, header)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message header")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("header saved"))
	return
}

// UploadMessageCatalog deals with uploading message catalogs
func UploadMessageCatalog(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	// Did we already upload the header?
	if message.IncomingPathExists(msgID, "catalog") {
		ErrorOut(w, http.StatusConflict, "catalog already uploaded")
		return
	}

	err := message.StoreCatalog(msgID, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message header")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved catalog"))
	return
}

// UploadMessageBlock deals with uploading message blocks
func UploadMessageBlock(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]
	id := vars["id"]

	// Did we already upload the header?
	if message.IncomingPathExists(msgID, id) {
		ErrorOut(w, http.StatusConflict, "block already uploaded")
		return
	}

	err := message.StoreBlock(msgID, id, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message block")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved message block"))
	return
}

// CompleteUpload is called whenever everything from a message has been uploaded and can be actually send
func CompleteUpload(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	// @TODO: How do we know if all data is send over? Maybe we should add all files to the header so we can verify?
	if !message.IncomingPathExists(msgID, "header.json") || !message.IncomingPathExists(msgID, "catalog") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	// queue the message for processing
	processor.QueueIncomingMessage(msgID)
}

// DeleteUpload is called whenever we want to completely remove a message by user request
func DeleteUpload(w http.ResponseWriter, req *http.Request) {
	// Delete the message and contents
	vars := mux.Vars(req)
	msgID := vars["msgid"]

	if !message.IncomingPathExists(msgID, "") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	err := message.RemoveMessage(message.SectionIncoming, msgID)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while deleting outgoing message")
		return
	}

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
