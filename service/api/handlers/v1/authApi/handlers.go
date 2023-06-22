package authApi

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/service"
	"github.com/ole-larsen/uploader/service/settings"
)

type API struct {
	service *service.Service
}

func NewAuthAPI(s *service.Service) *API {
	return &API{service: s}
}

func (a *API) XTokenAuth(token string) (*models.Principal, error) {
	if token == settings.Settings.XToken {
		principal := models.Principal(token)
		return &principal, nil
	}
	return nil, errors.New(http.StatusUnauthorized, "incorrect api key authApi")
}
