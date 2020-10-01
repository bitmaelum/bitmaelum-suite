package resolver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/ernesto-jimenez/httplogger"
	"github.com/sirupsen/logrus"
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
	Serial    uint64          `json:"serial_number"`
}

// RoutingDownload is a JSON structure we download from a resolver server
type RoutingDownload struct {
	Hash      string          `json:"hash"`
	PublicKey bmcrypto.PubKey `json:"public_key"`
	Routing   string          `json:"routing"`
	Serial    uint64          `json:"serial_number"`
}

// OrganisationDownload is a JSON structure we download from a resolver server
type OrganisationDownload struct {
	Hash        string                        `json:"hash"`
	PublicKey   bmcrypto.PubKey               `json:"public_key"`
	Validations []organisation.ValidationType `json:"validations"`
	Serial      uint64                        `json:"serial_number"`
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
	kd, err := r.fetchAddress(addr)
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
	kd, err := r.fetchRouting(routingID)
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
	kd, err := r.fetchOrganisation(orgHash)
	if err != nil {
		return nil, err
	}

	return &OrganisationInfo{
		Hash:        kd.Hash,
		PublicKey:   kd.PublicKey,
		Validations: kd.Validations,
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
	// Do a prefetch so we can get the current serial number
	addr, err := address.NewHashFromHash(info.Hash)
	if err != nil {
		return err
	}
	kd, err := r.fetchAddress(*addr)
	if err != nil {
		return err
	}

	data := &map[string]string{
		"public_key": info.PublicKey.String(),
		"routing_id": info.RoutingID,
		"proof":      proof.String(),
	}

	url := r.BaseURL + "/address/" + info.Hash
	return r.upload(url, data, generateAddressSignature(info, privKey, kd.Serial))
}

func (r *remoteRepo) UploadRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error {
	// Do a prefetch so we can get the current serial number
	kd, err := r.fetchRouting(info.Hash)
	if err != nil {
		return err
	}

	data := &map[string]string{
		"public_key": info.PublicKey.String(),
		"routing":    info.Routing,
	}

	url := r.BaseURL + "/routing/" + info.Hash
	return r.upload(url, data, generateRoutingSignature(info, privKey, kd.Serial))
}

func (r *remoteRepo) UploadOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey, proof proofofwork.ProofOfWork) error {
	// Do a prefetch so we can get the current serial number
	org, err := address.NewOrgHash(info.Hash)
	if err != nil {
		return err
	}
	kd, err := r.fetchOrganisation(*org)
	if err != nil {
		return err
	}

	data := &map[string]interface{}{
		"public_key":  info.PublicKey.String(),
		"proof":       proof.String(),
		"validations": info.Validations,
	}

	url := r.BaseURL + "/organisation/" + info.Hash
	return r.upload(url, data, generateOrganisationSignature(info, privKey, kd.Serial))
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
	// Do a prefetch so we can get the current serial number
	addr, err := address.NewHashFromHash(info.Hash)
	if err != nil {
		return err
	}
	kd, err := r.fetchAddress(*addr)
	if err != nil {
		return err
	}

	url := r.BaseURL + "/address/" + info.Hash
	return r.delete(url, generateAddressSignature(info, privKey, kd.Serial))
}

func (r *remoteRepo) DeleteRouting(info *RoutingInfo, privKey bmcrypto.PrivKey) error {
	// Do a prefetch so we can get the current serial number
	kd, err := r.fetchRouting(info.Hash)
	if err != nil {
		return err
	}

	url := r.BaseURL + "/routing/" + info.Hash
	return r.delete(url, generateRoutingSignature(info, privKey, kd.Serial))
}

func (r *remoteRepo) DeleteOrganisation(info *OrganisationInfo, privKey bmcrypto.PrivKey) error {
	// Do a prefetch so we can get the current serial number
	org, err := address.NewOrgHash(info.Hash)
	if err != nil {
		return err
	}
	kd, err := r.fetchOrganisation(*org)
	if err != nil {
		return err
	}

	url := r.BaseURL + "/organisation/" + info.Hash
	return r.delete(url, generateOrganisationSignature(info, privKey, kd.Serial))
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

	switch v := v.(type) {
	case *http.Request:
		data, err = httputil.DumpRequest(v, true)
	case *http.Response:
		data, err = httputil.DumpResponse(v, true)
	}
	if err != nil {
		logrus.Tracef("%s\n\n", err)
		return
	}

	logrus.Tracef("%s\n\n", data)
}

func (r *remoteRepo) fetchAddress(addr address.HashAddress) (*AddressDownload, error) {
	url := r.BaseURL + "/address/" + addr.String()

	kd := &AddressDownload{}
	err := r.resolve(url, &kd)
	if err != nil {
		return nil, err
	}

	return kd, nil
}

func (r *remoteRepo) fetchRouting(routingID string) (*RoutingDownload, error) {
	url := r.BaseURL + "/routing/" + routingID

	rd := &RoutingDownload{}
	err := r.resolve(url, &rd)
	if err != nil {
		return nil, err
	}

	return rd, nil
}

func (r *remoteRepo) fetchOrganisation(addr address.HashOrganisation) (*OrganisationDownload, error) {
	url := r.BaseURL + "/organisation/" + addr.String()

	od := &OrganisationDownload{}
	err := r.resolve(url, &od)
	if err != nil {
		return nil, err
	}

	return od, nil
}
