package resolve

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/ernesto-jimenez/httplogger"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

type remoteRepo struct {
	BaseURL string
	client  *http.Client
}

// KeyUpload is a JSON structure we upload to a resolver server
type KeyUpload struct {
	PublicKey bmcrypto.PubKey `json:"public_key"`
	Address   string          `json:"address"`
	Pow       string          `json:"pow"`
}

// KeyDownload is a JSON structure we download from a resolver server
type KeyDownload struct {
	Hash      string          `json:"hash"`
	PublicKey bmcrypto.PubKey `json:"public_key"`
	Address   string          `json:"address"`
}

// NewRemoteRepository creates new remote resolve repository
func NewRemoteRepository(baseURL string, debug bool) Repository {
	var transport http.RoundTripper = &http.Transport{}

	if debug {
		// Wrap transport in debug logging
		transport = httplogger.NewLoggedTransport(transport, internal.NewHTTPLogger())
	}

	return &remoteRepo{
		BaseURL: baseURL,
		client: &http.Client{
			Transport: transport,
		},
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

func (r *remoteRepo) Upload(info *Info, privKey bmcrypto.PrivKey, pow proofofwork.ProofOfWork) error {
	data := &KeyUpload{
		PublicKey: info.PublicKey,
		Address:   info.Server,
		Pow:       pow.String(),
	}

	byteBuf, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", r.BaseURL+"/"+info.Hash, bytes.NewBuffer(byteBuf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+generateSignature(info, privKey))

	logHTTP(req, nil)
	response, err := r.client.Do(req)
	if err != nil {
		logHTTP(response, nil)
		return err
	}

	logHTTP(response, nil)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	return errors.New(string(body))
}

func (r *remoteRepo) Delete(info *Info, privKey bmcrypto.PrivKey) error {
	req, err := http.NewRequest("DELETE", r.BaseURL+"/"+info.Hash, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+generateSignature(info, privKey))
	logHTTP(req, nil)
	response, err := r.client.Do(req)
	if err != nil {
		logHTTP(response, err)
		return err
	}

	logHTTP(response, err)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil
	}

	return errors.New(string(body))
}

func logHTTP(v interface{}, err error) {
	if err != nil {
		logrus.Tracef("%s\n\n", err)
		return
	}

	var data []byte

	switch v.(type) {
	case *http.Request:
		data, err = httputil.DumpRequest(v.(*http.Request), true)
	case *http.Response:
		data, err = httputil.DumpResponse(v.(*http.Response), true)
	}
	if err != nil {
		logrus.Tracef("%s\n\n", err)
		return
	}

	logrus.Tracef("%s\n\n", data)
}
