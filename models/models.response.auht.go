package models

import "encoding/json"

func UnmarshalPiholeAuthResponse(data []byte) (*PiholeAuthResponse, error) {
	var r PiholeAuthResponse
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	if err := validate.Struct(r); err != nil {
		return nil, err
	}
	return &r, err
}

func (r *PiholeAuthResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type PiholeAuthResponse struct {
	Session PiholeSession `json:"session" validate:"required"`
	Took    float64       `json:"took"`
}

type PiholeSession struct {
	Valid    bool   `json:"valid"  validate:"required"`
	Totp     bool   `json:"totp"  validate:"required"`
	Sid      string `json:"sid"  validate:"required"`
	CSRF     string `json:"csrf"  validate:"required"`
	Validity int64  `json:"validity"  validate:"required"`
}
