// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PublicFile public file
//
// swagger:model publicFile
type PublicFile struct {

	// attributes
	Attributes *PublicFileAttributes `json:"attributes,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`
}

// Validate validates this public file
func (m *PublicFile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAttributes(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicFile) validateAttributes(formats strfmt.Registry) error {
	if swag.IsZero(m.Attributes) { // not required
		return nil
	}

	if m.Attributes != nil {
		if err := m.Attributes.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attributes")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("attributes")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this public file based on the context it is used
func (m *PublicFile) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAttributes(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicFile) contextValidateAttributes(ctx context.Context, formats strfmt.Registry) error {

	if m.Attributes != nil {

		if swag.IsZero(m.Attributes) { // not required
			return nil
		}

		if err := m.Attributes.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("attributes")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("attributes")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PublicFile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PublicFile) UnmarshalBinary(b []byte) error {
	var res PublicFile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// PublicFileAttributes public file attributes
//
// swagger:model PublicFileAttributes
type PublicFileAttributes struct {

	// alt
	Alt string `json:"alt,omitempty"`

	// caption
	Caption string `json:"caption,omitempty"`

	// created
	// Format: date
	Created strfmt.Date `json:"created,omitempty"`

	// created by id
	CreatedByID int64 `json:"created_by_id,omitempty"`

	// deleted
	// Format: date
	Deleted strfmt.Date `json:"deleted,omitempty"`

	// ext
	Ext string `json:"ext,omitempty"`

	// hash
	Hash string `json:"hash,omitempty"`

	// height
	Height int64 `json:"height,omitempty"`

	// mime
	Mime string `json:"mime,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// provider
	Provider string `json:"provider,omitempty"`

	// size
	Size float64 `json:"size,omitempty"`

	// thumb
	Thumb string `json:"thumb,omitempty"`

	// type
	Type string `json:"type,omitempty"`

	// updated
	// Format: date
	Updated strfmt.Date `json:"updated,omitempty"`

	// updated by id
	UpdatedByID int64 `json:"updated_by_id,omitempty"`

	// url
	URL string `json:"url,omitempty"`

	// width
	Width int64 `json:"width,omitempty"`
}

// Validate validates this public file attributes
func (m *PublicFileAttributes) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreated(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDeleted(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdated(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PublicFileAttributes) validateCreated(formats strfmt.Registry) error {
	if swag.IsZero(m.Created) { // not required
		return nil
	}

	if err := validate.FormatOf("attributes"+"."+"created", "body", "date", m.Created.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PublicFileAttributes) validateDeleted(formats strfmt.Registry) error {
	if swag.IsZero(m.Deleted) { // not required
		return nil
	}

	if err := validate.FormatOf("attributes"+"."+"deleted", "body", "date", m.Deleted.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *PublicFileAttributes) validateUpdated(formats strfmt.Registry) error {
	if swag.IsZero(m.Updated) { // not required
		return nil
	}

	if err := validate.FormatOf("attributes"+"."+"updated", "body", "date", m.Updated.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this public file attributes based on context it is used
func (m *PublicFileAttributes) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PublicFileAttributes) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PublicFileAttributes) UnmarshalBinary(b []byte) error {
	var res PublicFileAttributes
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
