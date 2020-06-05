package handler

import (
    "github.com/gorilla/mux"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/container"
    "github.com/jaytaph/mailv2/core/messagebox"
    "net/http"
    "strconv"
)

// Create account handler
func RetrieveBoxes(w http.ResponseWriter, req *http.Request) {
    addr := core.StringToHash(mux.Vars(req)["addr"])

    // Retrieve all boxes
    as := container.GetAccountService()
    boxes := as.FetchMessageBoxes(addr, "*")

    type MailBoxListOutput struct {
        Address string                  `json:"address"`
        Boxes   []messagebox.MailBoxInfo   `json:"boxes"`
    }

    output := &MailBoxListOutput{
        Address: addr.String(),
        Boxes:   boxes,
    }

    _ = JsonOut(w, output)
}

func RetrieveBox(w http.ResponseWriter, req *http.Request) {
    addr := core.StringToHash(mux.Vars(req)["addr"])

    name := mux.Vars(req)["box"]

    offset := getQueryInt(req, "offset", 0)
    limit := getQueryInt(req, "limit", 1000)

    as := container.GetAccountService()
    mail := as.FetchMailFromBox(addr, name, offset, limit)

    _ = JsonOut(w, mail)
}

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

func getQueryString(req *http.Request, key string, def string) string {
    q := req.URL.Query()

    v := q.Get(key)
    if v == "" {
        return def
    }

    return v
}
