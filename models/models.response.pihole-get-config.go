package models

import "encoding/json"

func UnmarshalPiholeGetConfigRespone(data []byte) (*PiholeGetConfigRespone, error) {
	var r PiholeGetConfigRespone
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	if err := validate.Struct(r); err != nil {
		return nil, err
	}

	return &r, err
}

func (r *PiholeGetConfigRespone) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PiholeGetConfigRespone struct {
	Config PiholeConfig `json:"config"  validate:"required"`
}

type PiholeConfig struct {
	DNS PiholeConfigDNS `json:"dns" validate:"required"`
}

type PiholeConfigDNS struct {
	CnameRecords []string `json:"cnameRecords"  validate:"required"`
}
