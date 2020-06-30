package handler

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// UploadMessageHeader deals with uploading message headers
func UploadMessageHeader(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	p, err := getOutgoingPath(uuid, "header.json")
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "header path not found")
		return
	}

	// @TODO: When we send message to multiple users, we have multiple headers.
	// Check if header is already present
	_, err = os.Stat(p)
	if err == nil {
		ErrorOut(w, http.StatusConflict, "header already uploaded")
		return
	}

	// Create path first
	err = os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot create path")
		return
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "cannot read body")
		return
	}

	header := &message.Header{}
	err = json.Unmarshal(data, &header)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect header format")
		return
	}

	// Copy body straight to catalog file
	headerFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot save header")
		return
	}

	headerLen, err := headerFile.Write(data)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = os.Remove(p)
		ErrorOut(w, http.StatusInternalServerError, "cannot save header")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk(fmt.Sprintf("saved header of %d bytes", headerLen)))
	return
}

// UploadMessageCatalog deals with uploading message catalogs
func UploadMessageCatalog(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	p, err := getOutgoingPath(uuid, "catalog")
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "catalog path not found")
		return
	}

	// Check if catalog is already present
	_, err = os.Stat(p)
	if err == nil {
		ErrorOut(w, http.StatusConflict, "catalog already uploaded")
		return
	}

	// Create path first
	err = os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot create path")
		return
	}

	// Copy body straight to catalog file
	catFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot save catalog")
		return
	}

	blockLen, err := io.Copy(catFile, req.Body)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = os.Remove(p)
		ErrorOut(w, http.StatusInternalServerError, "cannot save catalog")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk(fmt.Sprintf("saved catalog of %d bytes", blockLen)))
	return
}

// UploadMessageBlock deals with uploading message blocks
func UploadMessageBlock(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	id := vars["id"]

	p, err := getOutgoingPath(uuid, id)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "path not found")
		return
	}

	// Check if block is already present
	_, err = os.Stat(p)
	if err == nil {
		ErrorOut(w, http.StatusConflict, "block already uploaded")
		return
	}

	// Create path if needed
	err = os.MkdirAll(path.Dir(p), 0777)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot create path")
		return
	}

	// Copy body straight to block file
	blockFile, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot save block")
		return
	}

	blockLen, err := io.Copy(blockFile, req.Body)
	if err != nil {
		// Something went wrong, remove the file just in case something was already written
		_ = os.Remove(p)
		ErrorOut(w, http.StatusInternalServerError, "cannot save block")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk(fmt.Sprintf("saved block of %d bytes", blockLen)))
	return
}

// SendMessage is called whenever everything from a message has been uploaded and can be actually send
func SendMessage(w http.ResponseWriter, req *http.Request) {
	// Check how we want to send the message.
	// Current mode: right now
	// Or maybe: delayed
}

// DeleteMessage is called whenever we want to completely remove a message by user request
func DeleteMessage(w http.ResponseWriter, req *http.Request) {
	// Delete the message and contents
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	p, _ := getOutgoingPath(uuid, "")
	_, err := os.Stat(p)
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "message not found")
		return
	}

	err = os.RemoveAll(p)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "error while deleting outgoing message")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(StatusOk("message removed"))
	return
}

func getOutgoingPath(uuid string, file string) (string, error) {
	return homedir.Expand(path.Join(".outgoing", uuid, file))
}
