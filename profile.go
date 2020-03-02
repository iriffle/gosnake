package main

import (
	"encoding/json"
)

type Profile struct {
	Name  string
	Color string
}

func NewProfile(name, color string) *Profile {
	p := Profile{
		Name:  name,
		Color: color,
	}
	return &p
}

func (p *Profile) Encode() ([]byte, error) {
	var jsonData []byte
	jsonData, err := json.Marshal(p)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func (p *Profile) Decode(jsonData []byte) error {
	err := json.Unmarshal(jsonData, p)
	if err != nil {
		return err
	}
	return nil
}
