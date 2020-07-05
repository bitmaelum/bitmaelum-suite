package handler

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"net/http"
	"strings"
)

// HomePage Information header on root /
func HomePage(w http.ResponseWriter, req *http.Request) {
	// Simple enough so things like curl work
	htmlVersion := false
	if strings.Contains(req.Header.Get("Accept"), "text/html") {
		htmlVersion = true
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	logo := core.GetMonochromeASCIILogo()
	if htmlVersion {
		logo = strings.Replace(logo, "\n", "<br>", -1)
		_, _ = w.Write([]byte("<pre>" + logo + "</pre>"))
	} else {
		_, _ = w.Write([]byte(logo))
	}
}
