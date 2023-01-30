package model

import "errors"

type Movie struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (m *Movie) Validate() error {
	if m.Name == "" {
		return errors.New("Name is missing")
	}
	return nil
}
