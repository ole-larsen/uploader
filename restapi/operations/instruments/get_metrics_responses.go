// Code generated by go-swagger; DO NOT EDIT.

package instruments

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/ole-larsen/uploader/models"
)

// GetMetricsOKCode is the HTTP code returned for type GetMetricsOK
const GetMetricsOKCode int = 200

/*
GetMetricsOK ok

swagger:response getMetricsOK
*/
type GetMetricsOK struct {

	/*
	  In: Body
	*/
	Payload models.Metrics `json:"body,omitempty"`
}

// NewGetMetricsOK creates GetMetricsOK with default headers values
func NewGetMetricsOK() *GetMetricsOK {

	return &GetMetricsOK{}
}

// WithPayload adds the payload to the get metrics o k response
func (o *GetMetricsOK) WithPayload(payload models.Metrics) *GetMetricsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get metrics o k response
func (o *GetMetricsOK) SetPayload(payload models.Metrics) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetMetricsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty map
		payload = models.Metrics{}
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}
