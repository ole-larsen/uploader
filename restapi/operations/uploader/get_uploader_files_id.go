// Code generated by go-swagger; DO NOT EDIT.

package uploader

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"

	"github.com/ole-larsen/uploader/models"
)

// GetUploaderFilesIDHandlerFunc turns a function with the right signature into a get uploader files ID handler
type GetUploaderFilesIDHandlerFunc func(GetUploaderFilesIDParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetUploaderFilesIDHandlerFunc) Handle(params GetUploaderFilesIDParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetUploaderFilesIDHandler interface for that can handle valid get uploader files ID params
type GetUploaderFilesIDHandler interface {
	Handle(GetUploaderFilesIDParams, *models.Principal) middleware.Responder
}

// NewGetUploaderFilesID creates a new http.Handler for the get uploader files ID operation
func NewGetUploaderFilesID(ctx *middleware.Context, handler GetUploaderFilesIDHandler) *GetUploaderFilesID {
	return &GetUploaderFilesID{Context: ctx, Handler: handler}
}

/*
	GetUploaderFilesID swagger:route GET /uploader/files/{id} uploader getUploaderFilesId

get raw file version
*/
type GetUploaderFilesID struct {
	Context *middleware.Context
	Handler GetUploaderFilesIDHandler
}

func (o *GetUploaderFilesID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetUploaderFilesIDParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
