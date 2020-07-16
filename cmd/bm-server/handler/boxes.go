package handler

import (
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// RetrieveBoxes retrieves all message boxes for the given account
func RetrieveBoxes(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		// @TODO: Return error
		return
	}

	// Retrieve all boxes
	as := container.GetAccountService()
	boxes := as.FetchMessageBoxes(*haddr, "*")

	type MailBoxListOutput struct {
		Address string                `json:"address"`
		Boxes   []message.MailBoxInfo `json:"boxes"`
	}

	output := &MailBoxListOutput{
		Address: haddr.String(),
		Boxes:   boxes,
	}

	_ = JSONOut(w, output)
}

// RetrieveBox retrieves info about the given mailbox
func RetrieveBox(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.NewHashFromHash(mux.Vars(req)["addr"])
	if err != nil {
		// @TODO: Return error
		return
	}
	name := mux.Vars(req)["box"]

	offset := getQueryInt(req, "offset", 0)
	limit := getQueryInt(req, "limit", 1000)

	as := container.GetAccountService()
	mail := as.FetchListFromBox(*haddr, name, offset, limit)

	_ = JSONOut(w, mail)
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

// Returns the given query key as string, or returns the default value
func getQueryString(req *http.Request, key string, def string) string {
	q := req.URL.Query()

	v := q.Get(key)
	if v == "" {
		return def
	}

	return v
}
