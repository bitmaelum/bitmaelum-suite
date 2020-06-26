package handler

import (
	"net/http"
)

// Handler when a message header is posted
func SendMessage(w http.ResponseWriter, req *http.Request) {

	//ret := OutputHeaderType{
	//    Error: false,
	//    Status: BODY_ACCEPT,
	//    Description: "Accepting body for this header",
	//    BodyAccept: &BodyAcceptType{
	//        Path: "/incoming/" + path,
	//        Timeout: to.Format(time.RFC3339),
	//    },
	//}
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//_ = json.NewEncoder(w).Encode(ret)
}
