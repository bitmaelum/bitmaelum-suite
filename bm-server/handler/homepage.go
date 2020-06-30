package handler

import (
	"github.com/bitmaelum/bitmaelum-server/core"
	"net/http"
	"strings"
)

// HomePage Information header on root /
func HomePage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	logo := core.GetMonochromeASCIILogo()
	logo = strings.Replace(logo, "\n", "<br>", -1)
	_, _ = w.Write([]byte("<pre>" + logo + "</pre>"))
}
