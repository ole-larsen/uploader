package v1

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/restapi/operations/instruments"
	"github.com/ole-larsen/uploader/restapi/operations/public"
	"github.com/ole-larsen/uploader/restapi/operations/uploader"
)

type ApiV1Handler interface {
	XTokenAuth(token string) (*models.Principal, error)

	PublicGetPing(params public.GetPingParams, principal *models.Principal) middleware.Responder
	PublicGetMetrics(params instruments.GetMetricsParams) middleware.Responder
	PublicGetFiles(params public.GetFilesFileParams) middleware.Responder

	UploaderGetFiles(params uploader.GetUploaderFilesParams, principal *models.Principal) middleware.Responder
	UploaderGetFilesID(params uploader.GetUploaderFilesIDParams, principal *models.Principal) middleware.Responder
	UploaderPostFiles(params uploader.PostUploaderFilesParams, principal *models.Principal) middleware.Responder
	UploaderPutFiles(params uploader.PutUploaderFilesParams, principal *models.Principal) middleware.Responder
}
