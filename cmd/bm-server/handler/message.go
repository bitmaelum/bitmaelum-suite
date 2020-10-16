// Copyright (c) 2020 BitMaelum Authors
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
	"strconv"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	incorrectBox string = "incorrect box"
)

// GetMessage will return a message header and catalog
func GetMessage(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, incorrectBox)
		return
	}

	messageID := mux.Vars(req)["message"]

	ar := container.GetAccountRepo()
	header, _ := ar.FetchMessageHeader(*haddr, box, messageID)
	catalog, _ := ar.FetchMessageCatalog(*haddr, box, messageID)

	output := &api.Message{
		ID:      messageID,
		Header:  *header,
		Catalog: catalog,
	}

	_ = JSONOut(w, output)
}

// GetMessageBlock will return a message block
func GetMessageBlock(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, incorrectBox)
		return
	}

	messageID := mux.Vars(req)["message"]
	blockID := mux.Vars(req)["block"]

	ar := container.GetAccountRepo()
	block, err := ar.FetchMessageBlock(*haddr, box, messageID, blockID)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, incorrectBox)
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
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, incorrectBox)
		return
	}

	messageID := mux.Vars(req)["message"]
	attachmentID := mux.Vars(req)["attachment"]

	ar := container.GetAccountRepo()
	attachment, size, err := ar.FetchMessageAttachment(*haddr, box, messageID, attachmentID)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect attachment")
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
