// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Ping ping
//
// swagger:model ping
type Ping struct {

	// Result of method execution. `pong` in case of success
	// Required: true
	// Enum: ["pong"]
	Ping *string `json:"ping"`
}

// Validate validates this ping
func (m *Ping) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePing(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var pingTypePingPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["pong"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		pingTypePingPropEnum = append(pingTypePingPropEnum, v)
	}
}

const (

	// PingPingPong captures enum value "pong"
	PingPingPong string = "pong"
)

// prop value enum
func (m *Ping) validatePingEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, pingTypePingPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *Ping) validatePing(formats strfmt.Registry) error {

	if err := validate.Required("ping", "body", m.Ping); err != nil {
		return err
	}

	// value enum
	if err := m.validatePingEnum("ping", "body", *m.Ping); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this ping based on context it is used
func (m *Ping) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Ping) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Ping) UnmarshalBinary(b []byte) error {
	var res Ping
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
