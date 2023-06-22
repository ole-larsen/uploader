// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/restapi/operations"
	"github.com/ole-larsen/uploader/restapi/operations/instruments"
	"github.com/ole-larsen/uploader/restapi/operations/public"
	"github.com/ole-larsen/uploader/restapi/operations/uploader"
	"github.com/ole-larsen/uploader/service"
	v1 "github.com/ole-larsen/uploader/service/api/handlers/v1"
	"github.com/ole-larsen/uploader/service/settings"
	"github.com/ole-larsen/uploader/service/util"
)

//go:generate swagger generate server --target ../../uploader --name Uploader --spec ../schema/swagger.yml --principal models.Principal

func configureFlags(api *operations.UploaderAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.UploaderAPI) http.Handler {
	s := service.NewService()

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf
	api.Logger = s.Logger.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	// register handlers for api v1
	handlersv1 := v1.NewApiV1Handlers(s)

	// setup file server to serve images and files
	http.Handle("/", http.FileServer(http.Dir("./uploads")))

	// Applies when the "Authorization" header is set

	// Applies when the "x-token" header is set
	api.XTokenAuth = handlersv1.XTokenAuth

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// uploader.PostUploaderFilesMaxParseMemory = 32 << 20
	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// uploader.PutUploaderFilesMaxParseMemory = 32 << 20
	api.PublicGetPingHandler = public.GetPingHandlerFunc(handlersv1.PublicGetPing)

	api.InstrumentsGetMetricsHandler = instruments.GetMetricsHandlerFunc(handlersv1.PublicGetMetrics)

	api.PublicGetFilesFileHandler = public.GetFilesFileHandlerFunc(handlersv1.PublicGetFiles)

	s.Logger.Infoln(settings.Settings.UseDB)

	if settings.Settings.UseDB == true {
		api.UploaderGetUploaderFilesHandler = uploader.GetUploaderFilesHandlerFunc(handlersv1.UploaderGetFiles)
		api.UploaderGetUploaderFilesIDHandler = uploader.GetUploaderFilesIDHandlerFunc(handlersv1.UploaderGetFilesID)
		api.UploaderPutUploaderFilesHandler = uploader.PutUploaderFilesHandlerFunc(handlersv1.UploaderPutFiles)
	} else {
		api.UploaderGetUploaderFilesHandler = uploader.GetUploaderFilesHandlerFunc(func(params uploader.GetUploaderFilesParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation uploader.GetUploaderFiles has not yet been implemented for using without db")
		})

		api.UploaderGetUploaderFilesIDHandler = uploader.GetUploaderFilesIDHandlerFunc(func(params uploader.GetUploaderFilesIDParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation uploader.GetUploaderFilesID has not yet been implemented for using without db")
		})

		api.UploaderPutUploaderFilesHandler = uploader.PutUploaderFilesHandlerFunc(func(params uploader.PutUploaderFilesParams, principal *models.Principal) middleware.Responder {
			return middleware.NotImplemented("operation uploader.PutUploaderFiles has not yet been implemented for using without db")
		})
	}

	api.UploaderPostUploaderFilesHandler = uploader.PostUploaderFilesHandlerFunc(handlersv1.UploaderPostFiles)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	handler = util.SetupCsrfHandler(handler)
	handler = util.SetupPrometheusHandler(handler)
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handler = util.SetupCorsHandler(handler)
	return handler
}
