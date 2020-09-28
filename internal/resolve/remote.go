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

// AddressDownload is a JSON structure we download from a resolver server
type AddressDownload struct {
	Hash      string          `json:"hash"`
	PublicKey bmcrypto.PubKey `json:"public_key"`
	RoutingID string          `json:"routing_id"`
}

// RoutingDownload is a JSON structure we download from a resolver server
type RoutingDownload struct {
	Hash      string          `json:"hash"`
	PublicKey bmcrypto.PubKey `json:"public_key"`
	Routing   string          `json:"routing"`
}

// OrganisationDownload is a JSON structure we download from a resolver server
type OrganisationDownload struct {
	Hash      string          `json:"hash"`
	PublicKey bmcrypto.PubKey `json:"public_key"`
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
func (r *remoteRepo) ResolveAddress(addr address.HashAddress) (*AddressInfo, error) {
	url := r.BaseURL + "/address/" + addr.String()

	kd := &AddressDownload{}
	err := r.resolve(url, &kd)
	if err != nil {
		return nil, err
	}

	return &AddressInfo{
		Hash:      kd.Hash,
		PublicKey: kd.PublicKey,
		RoutingID: kd.RoutingID,
	}, nil
}

func (r *remoteRepo) ResolveRouting(routingID string) (*RoutingInfo, error) {
	url := r.BaseURL + "/routing/" + routingID

	kd := &RoutingDownload{}
	err := r.resolve(url, &kd)
	if err != nil {
		return nil, err
	}

	return &RoutingInfo{
		Hash:      kd.Hash,
		PublicKey: kd.PublicKey,
		Routing:   kd.Routing,
	}, nil
}

func (r *remoteRepo) ResolveOrganisation(orgHash address.HashOrganisation) (*OrganisationInfo, error) {
	url := r.BaseURL + "/org/" + orgHash.String()

	kd := &OrganisationDownload{}
	err := r.resolve(url, &kd)
	if err != nil {
		return nil, err
	}

	return &OrganisationInfo{
		Hash:      kd.Hash,
		PublicKey: kd.PublicKey,
	}, nil
}

func (r *remoteRepo) resolve(url string, v interface{}) error {
	response, err := r.client.Get(url)
	if err != nil {
		logrus.Debugf("cannot get response from remote resolver: %s", err)
		return errKeyNotFound
	}

	if response.StatusCode == 404 {
		return errKeyNotFound
	}

	if response.StatusCode == 200 {
		res, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logrus.Debugf("cannot get body response from remote resolver: %s", err)
			return errKeyNotFound
		}

		err = json.Unmarshal(res, v)
		if err != nil {
			logrus.Debugf("cannot unmarshal resolve body: %s", err)
			return errKeyNotFound
		}

		return nil
	}

	return errKeyNotFound
}

func (r *remoteRepo) UploadAddress(info *AddressInfo, privKey bmcrypto.PrivKey, proof proofofwork.ProofOfWork) error {
	data := &map[string]string{
		"public_key": info.PublicKey.String(),
		"routing_id": info.RoutingID,
		"proof":      proof.String(),
	}

	url := r.BaseURL + "/address/" + info.Hash
	return r.upload(url, data, generateAddressSignature(info, privKey))
}

func (r *remoteRepo) UploadRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error {
	data := &map[string]string{
		"public_key": info.PublicKey.String(),
		"routing":    info.Routing,
	}

	url := r.BaseURL + "/routing/" + info.Hash
	return r.upload(url, data, generateRoutingSignature(info, privKey))
}

func (r *remoteRepo) UploadOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey, proof proofofwork.ProofOfWork) error {
	data := &map[string]string{
		"public_key": info.PublicKey.String(),
		"proof":      proof.String(),
	}

	url := r.BaseURL + "/organisation/" + info.Hash
	return r.upload(url, data, generateOrganisationSignature(info, privKey))
}

func (r *remoteRepo) upload(url string, v interface{}, sig string) error {
	byteBuf, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteBuf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sig)

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

func (r *remoteRepo) DeleteAddress(info *AddressInfo, privKey bmcrypto.PrivKey) error {
	url := r.BaseURL + "/address/" + info.Hash
	return r.delete(url, generateAddressSignature(info, privKey))
}

func (r *remoteRepo) DeleteRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error {
	url := r.BaseURL + "/routing/" + info.Hash
	return r.delete(url, generateRoutingSignature(info, privKey))
}

func (r *remoteRepo) DeleteOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey) error {
	url := r.BaseURL + "/org/" + info.Hash
	return r.delete(url, generateOrganisationSignature(info, privKey))
}

func (r *remoteRepo) delete(url, sig string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+sig)
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
