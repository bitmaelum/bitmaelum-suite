package api

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/internal/message"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	"io"
)

// UploadHeader uploads a header
func (api *API) UploadHeader(addr address.HashAddress, messageID string, header *message.Header) error {
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/account/%s/send/%s/header", addr.String(), messageID)
	return api.PostBytes(url, data)
}

// UploadCatalog uploads a catalog
func (api *API) UploadCatalog(addr address.HashAddress, messageID string, encryptedCatalog []byte) error {
	url := fmt.Sprintf("/account/%s/send/%s/catalog", addr.String(), messageID)
	return api.PostBytes(url, encryptedCatalog)
}

// UploadBlock uploads a message block or attachment
func (api *API) UploadBlock(addr address.HashAddress, messageID, blockID string, r io.Reader) error {
	url := fmt.Sprintf("/account/%s/send/%s/block/%s", addr.String(), messageID, blockID)
	return api.PostReader(url, r)
}

// DeleteMessage deletes a message and all content
func (api *API) DeleteMessage(addr address.HashAddress, messageID string) error {
	url := fmt.Sprintf("/account/%s/send/%s", addr.String(), messageID)
	return api.Delete(url)
}

// UploadComplete signals the mailserver that all blocks (and headers) have been uploaded and can start processing
func (api *API) UploadComplete(addr address.HashAddress, messageID string) error {
	url := fmt.Sprintf("/account/%s/send/%s", addr.String(), messageID)
	return api.PostBytes(url, []byte{})
}
