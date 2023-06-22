package v1

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/restapi/operations/instruments"
	"github.com/ole-larsen/uploader/restapi/operations/public"
	"github.com/ole-larsen/uploader/restapi/operations/uploader"
	"github.com/ole-larsen/uploader/service"
	"github.com/ole-larsen/uploader/service/api/handlers/v1/authApi"
	"github.com/ole-larsen/uploader/service/api/handlers/v1/publicApi"
	"github.com/ole-larsen/uploader/service/api/handlers/v1/uploaderApi"
)

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte("\n" + header + ": \n"))
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return nil
}

type Handlers struct {
	auth     *authApi.API
	public   *publicApi.API
	uploader *uploaderApi.API
}

func NewApiV1Handlers(s *service.Service) ApiV1Handler {
	return &Handlers{
		auth:     authApi.NewAuthAPI(s),
		public:   publicApi.NewPublicAPI(s),
		uploader: uploaderApi.NewUploaderAPI(s),
	}
}

func (h *Handlers) XTokenAuth(token string) (*models.Principal, error) {
	return h.auth.XTokenAuth(token)
}

func (h *Handlers) PublicGetPing(params public.GetPingParams, principal *models.Principal) middleware.Responder {
	_ = dumpRequest(os.Stdout, "ping", params.HTTPRequest)
	return h.public.GetPing(params, principal)
}

func (h *Handlers) PublicGetMetrics(params instruments.GetMetricsParams) middleware.Responder {
	_ = dumpRequest(os.Stdout, "metrics", params.HTTPRequest)
	return h.public.GetMetrics(params)
}

func (h *Handlers) UploaderGetFilesID(params uploader.GetUploaderFilesIDParams, principal *models.Principal) middleware.Responder {
	_ = dumpRequest(os.Stdout, "get_file_", params.HTTPRequest) // Ignore the error
	return h.uploader.GetFilesID(params, principal)
}

func (h *Handlers) PublicGetFiles(params public.GetFilesFileParams) middleware.Responder {
	_ = dumpRequest(os.Stdout, "get_file_", params.HTTPRequest) // Ignore the error
	return h.public.GetFilesFile(params)
}

func (h *Handlers) UploaderGetFiles(params uploader.GetUploaderFilesParams, principal *models.Principal) middleware.Responder {
	_ = dumpRequest(os.Stdout, "get_files", params.HTTPRequest) // Ignore the error
	return h.uploader.GetFiles(params, principal)
}

func (h *Handlers) UploaderPostFiles(params uploader.PostUploaderFilesParams, principal *models.Principal) middleware.Responder {
	_ = dumpRequest(os.Stdout, "post_files", params.HTTPRequest) // Ignore the error
	return h.uploader.PostFiles(params, principal)
}

func (h *Handlers) UploaderPutFiles(params uploader.PutUploaderFilesParams, principal *models.Principal) middleware.Responder {
	_ = dumpRequest(os.Stdout, "put_files", params.HTTPRequest) // Ignore the error
	return h.uploader.PutFiles(params, principal)
}
