package publicApi

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/nfnt/resize"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/restapi/operations/instruments"
	"github.com/ole-larsen/uploader/restapi/operations/public"
	"github.com/ole-larsen/uploader/service"
	"github.com/ole-larsen/uploader/service/api/handlers/v1/uploaderApi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/image/draw"
)

type API struct {
	service *service.Service
}

func NewPublicAPI(s *service.Service) *API {
	return &API{service: s}
}

func (a *API) GetPing(params public.GetPingParams, principal *models.Principal) middleware.Responder {
	pong := "pong"
	response := &models.Ping{Ping: &pong}
	return public.NewGetPingOK().WithPayload(response)
}

func (a *API) GetMetrics(params instruments.GetMetricsParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		promhttp.Handler().ServeHTTP(w, params.HTTPRequest)
	})
}

func (a *API) GetFilesFile(params public.GetFilesFileParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		path := strings.Split(params.HTTPRequest.URL.RequestURI(), "/")

		encodedFilename := path[len(path)-1]

		filename, err := url.QueryUnescape(encodedFilename)
		if err != nil {
			a.service.Logger.Errorln(err)
		}
		dimensions := strings.Split(filename, ":")[0]
		a.service.Logger.Infoln("try to serve", path, filename, dimensions, filename != dimensions)

		// if filename == dimensions means there are no dimension sizes

		if filename != dimensions {
			// create folder dimensions
			if _, err = os.Stat(uploaderApi.UPLOAD_DIR + "/" + dimensions); os.IsNotExist(err) {
				err = os.MkdirAll(uploaderApi.UPLOAD_DIR+"/"+dimensions, os.ModePerm)
				if err != nil {
					a.service.Logger.Errorln(err)
				}
			}
			files, err := ioutil.ReadDir(uploaderApi.UPLOAD_DIR + "/" + dimensions)
			if err != nil {
				a.service.Logger.Errorln(err)
			}

			filename = strings.Split(filename, ":")[1]

			for _, f := range files {
				if f.Name() == fmt.Sprintf("%s", filename) {
					a.serveFile(w, uploaderApi.UPLOAD_DIR+"/"+dimensions, filename)
					return
				}
			}
		}
		if filename == dimensions {
			a.service.Logger.Infoln("serve file without resize")
			a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
			return
		}

		// resize picture and put it to folder
		a.service.Logger.Infoln("resize and create file", filename)
		// Create a new file in the uploads directory
		input, err := os.Open(uploaderApi.UPLOAD_DIR + "/" + filename)
		if err != nil {
			a.service.Logger.Errorln(err)
		}

		defer func(input *os.File) {
			err = input.Close()
			if err != nil {
				a.service.Logger.Errorln(err)
			}
		}(input)

		fileParts := strings.Split(input.Name(), ".")
		extension := fileParts[len(fileParts)-1]

		switch extension {
		case "png":
			err = a.decodePNG(input, filename, dimensions)
			if err != nil {
				a.service.Logger.Errorln(err)
			} else {
				a.serveFile(w, uploaderApi.UPLOAD_DIR+"/"+dimensions, filename)
			}
			break
		case "jpg":
			err = a.decodeJPG(input, filename, dimensions)
			if err != nil {
				a.service.Logger.Errorln(err)
			} else {
				a.serveFile(w, uploaderApi.UPLOAD_DIR+"/"+dimensions, filename)
			}
			break
		case "pdf":
			a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
			break
		default:
			a.service.Logger.Fatalln(extension)
			break
		}
	})
}

func (a *API) decodePNG(input *os.File, filename string, dimensions string) error {
	width, err := strconv.Atoi(strings.Split(dimensions, "x")[0])
	if err != nil {
		return err
	}
	height, err := strconv.Atoi(strings.Split(dimensions, "x")[1])
	if err != nil {
		return err
	}

	src, err := png.Decode(input)
	if err != nil {
		return err
	}
	// Set the expected size that you want:
	resized := image.NewRGBA(image.Rect(0, 0, width, height))

	// Resize:
	draw.NearestNeighbor.Scale(resized, resized.Rect, src, src.Bounds(), draw.Over, nil)

	dst, err := os.Create(fmt.Sprintf("%s/%s", uploaderApi.UPLOAD_DIR+"/"+dimensions, filename))
	if err != nil {
		return err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	return png.Encode(dst, resized)
}

func (a *API) decodeJPG(input *os.File, filename string, dimensions string) error {
	width, err := strconv.Atoi(strings.Split(dimensions, "x")[0])
	if err != nil {
		return err
	}
	height, err := strconv.Atoi(strings.Split(dimensions, "x")[1])
	if err != nil {
		return err
	}

	src, _, err := image.Decode(input)
	if err != nil {
		return err
	}

	resized := resize.Resize(uint(width), uint(height), src, resize.Lanczos3)

	dst, err := os.Create(fmt.Sprintf("%s/%s", uploaderApi.UPLOAD_DIR+"/"+dimensions, filename))
	if err != nil {
		return err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	return jpeg.Encode(dst, resized, nil)
}

func (a *API) serveFile(w http.ResponseWriter, path string, filename string) {
	buf, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, filename))

	if err != nil {
		a.service.Logger.Errorln(err)
	}

	ext := strings.Replace(filepath.Ext(filename), ".", "", 1)
	a.service.Logger.Infoln("serve", ext)
	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", ext))

	if ext == "svg" {
		w.Header().Set("Content-Type", fmt.Sprintf("image/%s+xml", ext))
	}

	if ext == "jpeg" {
		w.Header().Set("Content-Type", fmt.Sprintf("image/%s", "jpg"))
	}

	if ext == "pdf" {
		w.Header().Set("Content-type", "application/pdf")
	}

	w.Write(buf)
}
