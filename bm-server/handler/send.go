package handler

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-server/bm-server/foobar"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

// UploadMessageHeader deals with uploading message headers
func UploadMessageHeader(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	// Did we already upload the header?
	if message.IncomingPathExists(uuid, "header.json") {
		ErrorOut(w, http.StatusConflict, "header already uploaded")
		return
	}

	// Read header from request body
	header, err := readHeaderFromBody(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect header structure")
		return
	}

	// Save request
	err = message.StoreMessageHeader(uuid, header)
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
	uuid := vars["uuid"]

	// Did we already upload the header?
	if message.IncomingPathExists(uuid, "catalog") {
		ErrorOut(w, http.StatusConflict, "catalog already uploaded")
		return
	}

	err := message.StoreCatalog(uuid, req.Body)
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
	uuid := vars["uuid"]
	id := vars["id"]

	// Did we already upload the header?
	if message.IncomingPathExists(uuid, id) {
		ErrorOut(w, http.StatusConflict, "block already uploaded")
		return
	}

	err := message.StoreBlock(uuid, id, req.Body)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while storing message block")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("saved message block"))
	return
}

// SendMessage is called whenever everything from a message has been uploaded and can be actually send
func SendMessage(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	if ! message.IncomingPathExists(uuid, "") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	// Send uuid to client processor
	foobar.ProcessIncomingClientMessage(uuid)
}

// DeleteMessage is called whenever we want to completely remove a message by user request
func DeleteMessage(w http.ResponseWriter, req *http.Request) {
	// Delete the message and contents
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	if ! message.IncomingPathExists(uuid, "") {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	err := message.RemoveMessage(uuid)
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
