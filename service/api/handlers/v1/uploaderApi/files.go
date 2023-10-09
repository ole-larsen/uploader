package uploaderApi

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/restapi/operations/uploader"
)

func (a *API) GetFilesID(params uploader.GetUploaderFilesIDParams, principal *models.Principal) middleware.Responder {
	payload, err := a.getFilesID(params, principal)
	if err != nil {
		a.service.Logger.Errorln(err)
		code := int64(http.StatusInternalServerError)
		return uploader.NewGetUploaderFilesInternalServerError().WithPayload(&models.Error{
			Code:    code,
			Error:   true,
			Message: err.Error(),
		})
	}
	return uploader.NewGetUploaderFilesIDOK().WithPayload(payload)
}

func (a *API) GetFiles(params uploader.GetUploaderFilesParams, principal *models.Principal) middleware.Responder {
	payload, err := a.getFiles(params, principal)
	if err != nil {
		a.service.Logger.Errorln(err)
		code := int64(http.StatusInternalServerError)
		return uploader.NewGetUploaderFilesInternalServerError().WithPayload(&models.Error{
			Code:    code,
			Error:   true,
			Message: err.Error(),
		})
	}
	return uploader.NewGetUploaderFilesOK().WithPayload(payload)
}

func (a *API) PostFiles(params uploader.PostUploaderFilesParams, principal *models.Principal) middleware.Responder {
	payload, err := a.postFiles(params, principal)
	if err != nil {
		a.service.Logger.Errorln(err)
		code := int64(http.StatusInternalServerError)
		return uploader.NewPostUploaderFilesInternalServerError().WithPayload(&models.Error{
			Code:    code,
			Error:   true,
			Message: err.Error(),
		})
	}
	return uploader.NewPostUploaderFilesOK().WithPayload(payload)
}

func (a *API) PutFiles(params uploader.PutUploaderFilesParams, principal *models.Principal) middleware.Responder {
	payload, err := a.putFiles(params, principal)
	if err != nil {
		a.service.Logger.Errorln(err)
		code := int64(http.StatusInternalServerError)
		return uploader.NewGetUploaderFilesInternalServerError().WithPayload(&models.Error{
			Code:    code,
			Error:   true,
			Message: err.Error(),
		})
	}
	return uploader.NewPutUploaderFilesOK().WithPayload(payload)
}
