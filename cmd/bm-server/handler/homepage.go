package handler

import (
	"bytes"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
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

	logo := internal.GetMonochromeASCIILogo()
	if config.Server.Server.VerboseInfo {
		host := fmt.Sprintf("<<< %s >>>", config.Server.Server.Name)
		host = fmt.Sprintf("%*s ", (49+len(host))/2, host)
		logo := internal.GetMonochromeASCIILogo() + "\n\n" + host + "\n\n"

		var version bytes.Buffer
		internal.WriteVersionInfo("BitMealum-Server", &version)
		logo = logo + "\n\n" + version.String()
	}

	if htmlVersion {
		logo = strings.Replace(logo, "\n", "<br>", -1)
		_, _ = w.Write([]byte("<pre>" + logo + "</pre>"))
	} else {
		_, _ = w.Write([]byte(logo))
	}
}
