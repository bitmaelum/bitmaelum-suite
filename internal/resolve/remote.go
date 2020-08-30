package resolve

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type remoteRepo struct {
	BaseURL string
	client  *http.Client
}

// KeyUpload is a JSON structure we upload to a resolver server
type KeyUpload struct {
	PublicKey encrypt.PubKey `json:"public_key"`
	Address   string         `json:"address"`
	Signature string         `json:"signature"`
}

// KeyDownload is a JSON structure we download from a resolver server
type KeyDownload struct {
	Hash      string         `json:"hash"`
	PublicKey encrypt.PubKey `json:"public_key"`
	Address   string         `json:"address"`
}

// NewRemoteRepository creates new remote resolve repository
func NewRemoteRepository(baseURL string) Repository {
	return &remoteRepo{
		BaseURL: baseURL,
		client:  &http.Client{},
	}
}

// Resolve
func (r *remoteRepo) Resolve(addr address.HashAddress) (*Info, error) {
	response, err := r.client.Get(r.BaseURL + "/" + addr.String())
	if err != nil {
		logrus.Debugf("cannot get response from remote resolver: %s", err)
		return nil, errKeyNotFound
	}

	if response.StatusCode == 404 {
		return nil, errKeyNotFound
	}

	if response.StatusCode == 200 {
		res, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logrus.Debugf("cannot get body response from remote resolver: %s", err)
			return nil, errKeyNotFound
		}

		kd := &KeyDownload{}
		err = json.Unmarshal(res, &kd)
		if err != nil {
			logrus.Debugf("cannot unmarshal resolve body: %s", err)
			return nil, errKeyNotFound
		}

		ri := &Info{
			Hash:      kd.Hash,
			PublicKey: kd.PublicKey,
			Server:    kd.Address,
		}
		err = json.Unmarshal(res, &ri)
		if err != nil {
			logrus.Debugf("cannot unmarshal resolve body: %s", err)
			return nil, errKeyNotFound
		}

		return ri, nil
	}

	return nil, errKeyNotFound
}

func (r *remoteRepo) Upload(addr address.HashAddress, pubKey encrypt.PubKey, address, signature string) error {
	data := &KeyUpload{
		PublicKey: pubKey,
		Address:   address,
		Signature: signature,
	}

	byteBuf, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	response, err := r.client.Post(r.BaseURL+"/"+addr.String(), "application/json", bytes.NewBuffer(byteBuf))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	return errors.New(string(body))
}
