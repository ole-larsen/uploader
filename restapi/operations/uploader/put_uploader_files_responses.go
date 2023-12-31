// Code generated by go-swagger; DO NOT EDIT.

package uploader

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ole-larsen/uploader/models"
)

// PutUploaderFilesOKCode is the HTTP code returned for type PutUploaderFilesOK
const PutUploaderFilesOKCode int = 200

/*
PutUploaderFilesOK An array of files

swagger:response putUploaderFilesOK
*/
type PutUploaderFilesOK struct {

	/*
	  In: Body
	*/
	Payload []*models.File `json:"body,omitempty"`
}

// NewPutUploaderFilesOK creates PutUploaderFilesOK with default headers values
func NewPutUploaderFilesOK() *PutUploaderFilesOK {

	return &PutUploaderFilesOK{}
}

// WithPayload adds the payload to the put uploader files o k response
func (o *PutUploaderFilesOK) WithPayload(payload []*models.File) *PutUploaderFilesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put uploader files o k response
func (o *PutUploaderFilesOK) SetPayload(payload []*models.File) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutUploaderFilesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.File, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// PutUploaderFilesInternalServerErrorCode is the HTTP code returned for type PutUploaderFilesInternalServerError
const PutUploaderFilesInternalServerErrorCode int = 500

/*
PutUploaderFilesInternalServerError When some error occurs

swagger:response putUploaderFilesInternalServerError
*/
type PutUploaderFilesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutUploaderFilesInternalServerError creates PutUploaderFilesInternalServerError with default headers values
func NewPutUploaderFilesInternalServerError() *PutUploaderFilesInternalServerError {

	return &PutUploaderFilesInternalServerError{}
}

// WithPayload adds the payload to the put uploader files internal server error response
func (o *PutUploaderFilesInternalServerError) WithPayload(payload *models.Error) *PutUploaderFilesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put uploader files internal server error response
func (o *PutUploaderFilesInternalServerError) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutUploaderFilesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
