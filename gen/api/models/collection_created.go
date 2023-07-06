// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CollectionCreated collection created
//
// swagger:model CollectionCreated
type CollectionCreated struct {

	// collection id
	CollectionID string `json:"collection_id,omitempty"`
}

// Validate validates this collection created
func (m *CollectionCreated) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CollectionCreated) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CollectionCreated) UnmarshalBinary(b []byte) error {
	var res CollectionCreated
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}