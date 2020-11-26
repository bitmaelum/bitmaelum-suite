package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
)

// JwtErrorFunc is a generic error handling function that can be attached to an API client. This will automatically
// trigger whenever an error (https response >= 400) is found. In this case, it will only check for the token-time
// error, which is returned when the time of the client is off. 
func JwtErrorFunc(_ *http.Request, resp *http.Response) {
	if resp.StatusCode != 401 {
		return
	}

	// Read body
	b, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	err := api.GetErrorFromResponse(b)
	if err != nil && err.Error() == "token time not valid" {
		fmt.Println("The connection to the server was unauthenticated because of timing issues. It seems that your computer time is not")
		fmt.Println("set to the current time. This causes issues in communication with the BitMaelum server. Please update your time and")
		fmt.Println("try again.\n")
		os.Exit(1)
	}

	// Whoops.. not an error. Let's pretend nothing happened and create a new buffer so we can read the body again
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
}
