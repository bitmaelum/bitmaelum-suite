package handler

import (
	"net/http"
)

// GetMessage will return a message header and catalog
func GetMessage(w http.ResponseWriter, req *http.Request) {
	// box, err := strconv.Atoi(mux.Vars(req)["box"])
	// if err != nil {
	// 	ErrorOut(w, http.StatusBadRequest, "incorrect box")
	// 	return
	// }

	// messageID := mux.Vars(req)["message"]
	//
	// ar := container.GetAccountRepo()
	// flags, _ := ar.FetchMessageHeader(*haddr, box, messageID)
}

// GetMessageBlock will return a message block
func GetMessageBlock(w http.ResponseWriter, req *http.Request) {
}

// GetMessageAttachment will return a message attachment
func GetMessageAttachment(w http.ResponseWriter, req *http.Request) {
}
