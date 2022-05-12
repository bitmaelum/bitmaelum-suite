// Copyright (c) 2022 BitMaelum Authors
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
	"io"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	incorrectBlock string = "incorrect block ID"
)

type moveMessageInput struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// DeleteMessage will delete a message
func DeleteMessage(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]

	ar := container.Instance.GetAccountRepo()
	err = ar.RemoveMessage(*haddr, messageID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}

	_ = httputils.JSONOut(w, http.StatusOK, nil)
}

// GetMessage will return a message header and catalog
func GetMessage(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]

	ar := container.Instance.GetAccountRepo()
	header, err := ar.FetchMessageHeader(*haddr, messageID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}
	catalog, err := ar.FetchMessageCatalog(*haddr, messageID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}

	output := &api.Message{
		ID:      messageID,
		Header:  *header,
		Catalog: catalog,
	}

	_ = httputils.JSONOut(w, http.StatusOK, output)
}

// CopyMessage will copy a message to a box
func CopyMessage(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]

	var input moveMessageInput
	err = httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	ar := container.Instance.GetAccountRepo()
	err = ar.CopyMessage(*haddr, messageID, input.To)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}

	httputils.JSONOut(w, http.StatusOK, "")
}

// MoveMessage will move a message to a box
func MoveMessage(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]

	var input moveMessageInput
	err = httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	ar := container.Instance.GetAccountRepo()
	err = ar.MoveMessage(*haddr, messageID, input.From, input.To)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}

	httputils.JSONOut(w, http.StatusOK, "")
}

// RemoveMessageFromBox will remove a message from a box
func RemoveMessageFromBox(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]

	var input moveMessageInput
	err = httputils.DecodeBody(w, req.Body, &input)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	ar := container.Instance.GetAccountRepo()
	err = ar.RemoveFromBox(*haddr, input.From, messageID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}

	// Should the message be removed if it's not linked to any mailbox?
	boxes, err := ar.GetAllBoxes(*haddr)
	if err == nil {
		var found bool
		for box := range boxes {
			if ar.ExistsInBox(*haddr, box, messageID) {
				found = true
			}
		}
		if !found {
			ar.RemoveMessage(*haddr, messageID)
		}
	}

	httputils.JSONOut(w, http.StatusOK, "")
}

// GetMessageBlock will return a message block
func GetMessageBlock(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]
	blockID := mux.Vars(req)["block"]

	ar := container.Instance.GetAccountRepo()
	block, err := ar.FetchMessageBlock(*haddr, messageID, blockID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, incorrectBlock)
		return
	}

	w.Header().Set("Content-Type", "application/binary")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(block)
}

// GetMessageAttachment will return a message attachment
func GetMessageAttachment(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		httputils.ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	messageID := mux.Vars(req)["message"]
	attachmentID := mux.Vars(req)["attachment"]

	ar := container.Instance.GetAccountRepo()
	attachment, size, err := ar.FetchMessageAttachment(*haddr, messageID, attachmentID)
	if err != nil {
		httputils.ErrorOut(w, http.StatusBadRequest, "incorrect attachment")
		return
	}

	// Copy data from our reader to http writer (buffered)
	w.Header().Set("Content-Type", "application/binary")
	w.WriteHeader(http.StatusOK)
	bw, err := io.Copy(w, attachment)
	if err != nil {
		logrus.Errorf("Could only write %d out of %d bytes from attachment %s/%s", bw, size, messageID, attachmentID)
	}
}
