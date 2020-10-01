package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
)

// UploadHeader uploads a header
func (api *API) UploadHeader(t ticket.Ticket, header *message.Header) error {
	api.setTicketHeader(t)
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return err
	}

	resp, statusCode, err := api.Post("/incoming/header", data)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return getErrorFromResponse(resp)
	}

	return nil
}

// UploadCatalog uploads a catalog
func (api *API) UploadCatalog(t ticket.Ticket, encryptedCatalog []byte) error {
	api.setTicketHeader(t)
	_, statusCode, err := api.Post("/incoming/catalog", encryptedCatalog)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

// UploadBlock uploads a message block or attachment
func (api *API) UploadBlock(t ticket.Ticket, blockID string, r io.Reader) error {
	api.setTicketHeader(t)
	url := fmt.Sprintf("/incoming/block/%s", blockID)
	_, statusCode, err := api.PostReader(url, r)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

// DeleteMessage deletes a message and all content
func (api *API) DeleteMessage(t ticket.Ticket) error {
	api.setTicketHeader(t)

	_, statusCode, err := api.Delete("/incoming")
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}

// CompleteUpload signals the mailserver that all blocks (and headers) have been uploaded and can start processing
func (api *API) CompleteUpload(t ticket.Ticket) error {
	api.setTicketHeader(t)
	_, statusCode, err := api.Post("/incoming", []byte{})
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return errNoSuccess
	}

	return nil
}
