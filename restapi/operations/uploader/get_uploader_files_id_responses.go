// Code generated by go-swagger; DO NOT EDIT.

package uploader

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ole-larsen/uploader/models"
)

// GetUploaderFilesIDOKCode is the HTTP code returned for type GetUploaderFilesIDOK
const GetUploaderFilesIDOKCode int = 200

/*
GetUploaderFilesIDOK file

swagger:response getUploaderFilesIdOK
*/
type GetUploaderFilesIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.File `json:"body,omitempty"`
}

// NewGetUploaderFilesIDOK creates GetUploaderFilesIDOK with default headers values
func NewGetUploaderFilesIDOK() *GetUploaderFilesIDOK {

	return &GetUploaderFilesIDOK{}
}

// WithPayload adds the payload to the get uploader files Id o k response
func (o *GetUploaderFilesIDOK) WithPayload(payload *models.File) *GetUploaderFilesIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get uploader files Id o k response
func (o *GetUploaderFilesIDOK) SetPayload(payload *models.File) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUploaderFilesIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUploaderFilesIDInternalServerErrorCode is the HTTP code returned for type GetUploaderFilesIDInternalServerError
const GetUploaderFilesIDInternalServerErrorCode int = 500

/*
GetUploaderFilesIDInternalServerError When some error occurs

swagger:response getUploaderFilesIdInternalServerError
*/
type GetUploaderFilesIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetUploaderFilesIDInternalServerError creates GetUploaderFilesIDInternalServerError with default headers values
func NewGetUploaderFilesIDInternalServerError() *GetUploaderFilesIDInternalServerError {

	return &GetUploaderFilesIDInternalServerError{}
}

// WithPayload adds the payload to the get uploader files Id internal server error response
func (o *GetUploaderFilesIDInternalServerError) WithPayload(payload *models.Error) *GetUploaderFilesIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get uploader files Id internal server error response
func (o *GetUploaderFilesIDInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUploaderFilesIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
