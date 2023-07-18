package models

type Name struct {
	Id int64 `json:"id,omitempty"`

	Firstname string `json:"firstname,omitempty"`

	Lastname string `json:"lastname,omitempty"`

	Name string `json:"name,omitempty"`
}
