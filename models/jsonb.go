package models

import "encoding/json"

type JSONB string

func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == "" {
		return []byte("null"), nil
	}
	return []byte(j), nil
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*j = ""
		return nil
	}

	if !json.Valid(data) {
		return &json.SyntaxError{}
	}

	*j = JSONB(data)
	return nil
}
