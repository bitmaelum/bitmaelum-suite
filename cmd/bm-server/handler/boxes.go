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
	"net/http"
	"strconv"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
)

type jsonOut map[string]interface{}

type boxIn struct {
	ParentBoxID int `json:"parent_box_id"`
}

const (
	accountNotFound string = "account not found"
)

// CreateBox creates a new box under a specific parent box or 0 for a root box
func CreateBox(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	var input boxIn
	err = DecodeBody(w, req.Body, &input)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	ar := container.GetAccountRepo()
	err = ar.CreateBox(*haddr, input.ParentBoxID)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// DeleteBox deletes a given box with all messages (note: what about child boxes??)
func DeleteBox(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "box not found")
		return
	}

	ar := container.GetAccountRepo()
	err = ar.DeleteBox(*haddr, box)
	if err != nil {
		ErrorOut(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// RetrieveBoxes retrieves all message boxes for the given account
func RetrieveBoxes(w http.ResponseWriter, req *http.Request) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	// Retrieve all boxes
	ar := container.GetAccountRepo()
	boxes, err := ar.GetAllBoxes(*haddr)
	if err != nil {
		ErrorOut(w, http.StatusInternalServerError, "cannot read boxes")
		return
	}

	output := jsonOut{
		"meta": jsonOut{
			"total":    len(boxes),
			"returned": len(boxes),
		},
		"boxes": boxes,
	}

	_ = JSONOut(w, output)
}

// RetrieveMessagesFromBox retrieves info about the given mailbox
func RetrieveMessagesFromBox(w http.ResponseWriter, req *http.Request) {
	list, err := getMessageList(req)
	if err != nil {
		ErrorOut(w, http.StatusNotFound, accountNotFound)
		return
	}

	// if we authenticated by API-key, we only return the message IDs
	if IsAPIKeyAuthenticated(req) {
		var ids []string
		for _, msg := range list.Messages {
			ids = append(ids, msg.ID)
		}

		_ = JSONOut(w, jsonOut{
			"meta":        list.Meta,
			"message_ids": ids,
		})
		return
	}

	// Otherwise, return the whole list
	_ = JSONOut(w, list)
}

func getMessageList(req *http.Request) (*account.MessageList, *httpError) {
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		return nil, &httpError{
			err:        accountNotFound,
			StatusCode: http.StatusNotFound,
		}
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		return nil, &httpError{
			err:        "box not found",
			StatusCode: http.StatusNotFound,
		}
	}

	ar := container.GetAccountRepo()
	if !ar.ExistsBox(*haddr, box) {
		return nil, &httpError{
			err:        accountNotFound,
			StatusCode: http.StatusNotFound,
		}
	}

	since := getQueryInt(req, "since", 0)
	sinceTs := time.Unix(int64(since), 0)
	offset := getQueryInt(req, "offset", 0)
	limit := getQueryInt(req, "limit", 100)

	list, err := ar.FetchListFromBox(*haddr, box, sinceTs, offset, limit)
	if err != nil {
		return nil, &httpError{
			err:        err.Error(),
			StatusCode: http.StatusInternalServerError,
		}
	}

	return list, nil
}

// Returns the given query key as integer, or returns the default value
func getQueryInt(req *http.Request, key string, def int) int {
	q := req.URL.Query()

	v := q.Get(key)
	if v == "" {
		return def
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}

	return i
}

// // Returns the given query key as string, or returns the default value
// func getQueryString(req *http.Request, key string, def string) string {
// 	q := req.URL.Query()
//
// 	v := q.Get(key)
// 	if v == "" {
// 		return def
// 	}
//
// 	return v
// }
