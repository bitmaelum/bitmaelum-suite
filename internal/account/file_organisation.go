package account

import (
	"encoding/json"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

func (r *fileRepo) StoreOrganisationSettings(addr address.HashAddress, settings OrganisationSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	// And store
	return r.store(addr, organisationFile, data)
}

func (r *fileRepo) FetchOrganisationSettings(addr address.HashAddress) (*OrganisationSettings, error) {
	settings := &OrganisationSettings{}
	err := r.fetchJSON(addr, keysFile, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
