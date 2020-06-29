package api

import (
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"io"
)

func (api *Api) UploadHeader(addr core.HashAddress, messageId string, header *message.Header) error {
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		return err
	}

	return api.PostBytes("/account/"+addr.String()+"/send/"+messageId+"/header", data)
}

func (api *Api) UploadCatalog(addr core.HashAddress, messageId string, encryptedCatalog []byte) error {
	return api.PostBytes("/account/"+addr.String()+"/send/"+messageId+"/catalog", encryptedCatalog)
}

func (api *Api) UploadBlock(addr core.HashAddress, messageId, blockId string, r io.Reader) error {
	return api.PostReader("/account/"+addr.String()+"/send/"+messageId+"/block/"+blockId, r)
}

func (api *Api) DeleteMessage(addr core.HashAddress, messageId string) error {
	return api.Delete("/account/" + addr.String() + "/send/" + messageId)
}
