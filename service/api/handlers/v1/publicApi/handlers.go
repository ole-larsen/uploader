package publicApi

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
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
		path := strings.Split(params.HTTPRequest.URL.Path, "/")
		dir := strings.Replace(strings.Join(path[:len(path)-1], "/"), "/api/v1/files", uploaderApi.UPLOAD_DIR, 1)

		encodedFilename := path[len(path)-1]

		filename, err := url.QueryUnescape(encodedFilename)
		if err != nil {
			a.internalError(w, err)
			return
		}
		ext := a.extractExt(filename)
		name := a.extractName(filename)

		if params.Format != nil && *params.Format != "auto" {
			ext = *params.Format
		}

		// create folder dimensions
		if params.W != nil {
			sWidth := fmt.Sprintf("%d", int(*params.W))
			if _, err = os.Stat(dir + "/" + sWidth); os.IsNotExist(err) {
				err = os.MkdirAll(dir+"/"+sWidth, os.ModePerm)
				if err != nil {
					a.internalError(w, err)
					return
				}
			}
			files, err := ioutil.ReadDir(dir + "/" + sWidth)
			if err != nil {
				a.internalError(w, err)
				return
			}
			if len(files) > 0 {
				for _, f := range files {
					if f.Name() == filename {
						a.serveFile(w, dir+"/"+sWidth, name+"."+ext)
						return
					}
				}
			}
		}

		if params.W == nil {
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				a.internalError(w, err)
				return
			}
			if len(files) > 0 {
				for _, f := range files {
					if f.Name() == name+"."+ext {
						a.serveFile(w, dir, name+"."+ext)
						return
					}
				}
			}
		}

		// open source file
		input, err := os.Open(dir + "/" + filename)
		if err != nil {
			a.internalError(w, err)
			return
		}

		defer func(input *os.File) {
			err = input.Close()
			if err != nil {
				a.service.Logger.Errorln(err)
			}
		}(input)

		// decode source file
		src, _, err := image.Decode(input)
		if err != nil {
			a.internalError(w, err)
			return
		}
		bounds := src.Bounds()

		sourceWidth := bounds.Dx()
		sourceHeight := bounds.Dy()

		width := float64(sourceWidth)
		height := float64(sourceHeight)

		coef := height / width

		if params.W != nil {
			width = *params.W
			height = coef * width
		}

		dpr := float64(1)

		if params.Dpr != nil {
			dpr = *params.Dpr
			width = width * dpr
			height = height * dpr
		}

		switch ext {
		case "webp":
			if params.W == nil {
				err = a.decodeBaseWEBP(src, dir, name+".webp")
				if err != nil {
					a.internalError(w, err)
				} else {
					a.serveFile(w, dir, name+".webp")
				}
			}
			if params.W != nil {
				err = a.decodeWEBP(src, dir, name+".webp", int(width), int(height))
				if err != nil {
					a.internalError(w, err)
				} else {
					a.serveFile(w, fmt.Sprintf("%s/%d", dir, int(width)), name+".webp")
				}
			}
		case "png":
			if params.W == nil {
				err = a.decodeBasePNG(src, dir, name+".png")
				if err != nil {
					a.internalError(w, err)
				} else {
					a.serveFile(w, dir, name+".png")
				}

			}
			if params.W != nil {
				err = a.decodePNG(src, dir, name+".png", int(width), int(height))
				if err != nil {
					a.internalError(w, err)
				} else {
					a.serveFile(w, dir, name+".png")
				}
			}
		case "jpg":
			if params.W == nil {
				err = a.decodeBaseJPG(src, dir, name+".jpg")
				if err != nil {
					a.internalError(w, err)
				} else {
					a.serveFile(w, dir, name+".jpg")
				}
			}
			if params.W != nil {
				err = a.decodeJPG(src, dir, name+".jpg", int(width), int(height))
				if err != nil {
					a.internalError(w, err)
				} else {
					a.serveFile(w, dir, name+".jpg")
				}
			}
		case "pdf":
			a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
		default:
			a.service.Logger.Fatalln(name, ext)
		}
	})
}

func (a *API) extractName(filename string) string {
	parts := strings.Split(filename, ".")
	return strings.Join(parts[:len(parts)-1], "/")
}

func (a *API) extractExt(filename string) string {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}

func (a *API) decodeBaseWEBP(src image.Image, dir string, filename string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		return err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	// Encode the image in WebP format
	return webp.Encode(dst, src, nil)
}

func (a *API) decodeWEBP(src image.Image, dir string, filename string, width int, height int) error {
	// Set the expected size that you want:
	resized := image.NewRGBA(image.Rect(0, 0, width, height))

	// Resize:
	draw.NearestNeighbor.Scale(resized, resized.Rect, src, src.Bounds(), draw.Over, &draw.Options{})

	dst, err := os.Create(fmt.Sprintf("%s/%d/%s", dir, width, filename))
	if err != nil {
		return err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	// Encode the image in WebP format
	return webp.Encode(dst, resized, nil)
}

func (a *API) decodeBasePNG(src image.Image, dir string, filename string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		return err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	return png.Encode(dst, src)
}

func (a *API) decodePNG(src image.Image, dir string, filename string, width int, height int) error {
	// Set the expected size that you want:
	resized := image.NewRGBA(image.Rect(0, 0, width, height))

	// Resize:
	draw.NearestNeighbor.Scale(resized, resized.Rect, src, src.Bounds(), draw.Over, &draw.Options{})

	dst, err := os.Create(fmt.Sprintf("%s/%d/%s", dir, width, filename))
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

func (a *API) decodeBaseJPG(src image.Image, dir string, filename string) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s", dir, filename))
	if err != nil {
		return err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	return jpeg.Encode(dst, src, nil)
}

func (a *API) decodeJPG(src image.Image, dir string, filename string, width int, height int) error {
	resized := resize.Resize(uint(width), uint(height), src, resize.Lanczos3)

	dst, err := os.Create(fmt.Sprintf("%s/%d/%s", dir, width, filename))
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

	a.service.Logger.Infoln("serve", filename, ext)

	contentType := fmt.Sprintf("image/%s", ext)

	if ext == "svg" {
		contentType = fmt.Sprintf("image/%s+xml", ext)
	}

	if ext == "jpeg" {
		contentType = fmt.Sprintf("image/%s", "jpg")
	}

	if ext == "pdf" {
		contentType = "application/pdf"
	}

	if ext == "webp" {
		contentType = fmt.Sprintf("image/%s", ext)
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(buf)
}

func (a *API) internalError(w http.ResponseWriter, err error) {
	a.service.Logger.Errorln(err)
	code := int64(http.StatusInternalServerError)
	internalError, _ := json.Marshal(models.Error{
		Code:    code,
		Error:   true,
		Message: err.Error(),
	})
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(internalError)
}
