package api

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"io"
)

// UploadHeader uploads a header
func (api *API) UploadHeader(addr core.HashAddress, messageID string, header *message.Header) error {
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return err
	}

	return api.PostBytes("/account/"+addr.String()+"/send/"+messageID+"/header", data)
}

// UploadCatalog uploads a catalog
func (api *API) UploadCatalog(addr core.HashAddress, messageID string, encryptedCatalog []byte) error {
	return api.PostBytes("/account/"+addr.String()+"/send/"+messageID+"/catalog", encryptedCatalog)
}

// UploadBlock uploads a message block or attachment
func (api *API) UploadBlock(addr core.HashAddress, messageID, blockID string, r io.Reader) error {
	return api.PostReader("/account/"+addr.String()+"/send/"+messageID+"/block/"+blockID, r)
}

// DeleteMessage deletes a message and all content
func (api *API) DeleteMessage(addr core.HashAddress, messageID string) error {
	return api.Delete("/account/" + addr.String() + "/send/" + messageID)
}
