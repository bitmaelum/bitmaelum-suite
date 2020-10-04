package api

import (
	"fmt"
	"io"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

// Message is a standard structure that returns a message header + catalog
type Message struct {
	ID      string         `json:"id"`
	Header  message.Header `json:"h"`
	Catalog []byte         `json:"c"`
}

// GetMessage retrieves a message header + catalog from a message box
func (api *API) GetMessage(addr address.Hash, box, messageID string) (*Message, error) {
	in := &Message{}

	url := fmt.Sprintf("/account/%s/box/%s/message/%s", addr.String(), box, messageID)
	resp, statusCode, err := api.GetJSON(url, in)
	if err != nil {
		logrus.Trace(err)
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, getErrorFromResponse(resp)
	}

	return in, nil
}

// GetMessageBlock retrieves a message block
func (api *API) GetMessageBlock(addr address.Hash, box, messageID, blockID string) ([]byte, error) {
	body, statusCode, err := api.Get(fmt.Sprintf("/account/%s/box/%s/message/%s/block/%s", addr.String(), box, messageID, blockID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return body, nil
}

// GetMessageAttachment retrieves a message attachment reader
func (api *API) GetMessageAttachment(addr address.Hash, box, messageID, attachmentID string) (io.Reader, error) {
	r, statusCode, err := api.GetReader(fmt.Sprintf("/account/%s/box/%s/message/%s/attachment/%s", addr.String(), box, messageID, attachmentID))
	if err != nil {
		return nil, err
	}

	if statusCode < 200 || statusCode > 299 {
		return nil, errNoSuccess
	}

	return r, nil
}
