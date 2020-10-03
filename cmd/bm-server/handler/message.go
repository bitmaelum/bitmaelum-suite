package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetMessage will return a message header and catalog
func GetMessage(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.HashFromString(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
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
	haddr, err := address.HashFromString(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
		return
	}

	messageID := mux.Vars(req)["message"]
	blockID := mux.Vars(req)["block"]

	ar := container.GetAccountRepo()
	block, err := ar.FetchMessageBlock(*haddr, box, messageID, blockID)
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect block")
		return
	}

	w.Header().Set("Content-Type", "application/binary")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(block)
}

// GetMessageAttachment will return a message attachment
func GetMessageAttachment(w http.ResponseWriter, req *http.Request) {
	haddr, err := address.HashFromString(mux.Vars(req)["addr"])
	if err != nil {
		ErrorOut(w, http.StatusNotFound, "account not found")
		return
	}

	box, err := strconv.Atoi(mux.Vars(req)["box"])
	if err != nil {
		ErrorOut(w, http.StatusBadRequest, "incorrect box")
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
