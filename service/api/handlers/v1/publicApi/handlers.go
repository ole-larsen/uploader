package publicApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
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
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
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
		fmt.Println("-------------------------------------")
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

		if params.W == nil && params.Dpr == nil && params.Format == nil {
			a.service.Logger.Infoln("[serve original]", dir, name, ext)
			a.serveOriginal(w, dir, name, ext)
			return
		}

		src := a.getSource(w, dir, filename, ext)

		if src == nil {
			return
		}

		if strings.Contains(ext, "pdf") {
			a.service.Logger.Infoln("[serve pdf]", dir, name, ext)
			a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
			return
		}

		width, height := a.getSize(src, params.W, params.Dpr)
		// create folder by dimensions if not exists
		if params.W != nil {
			sWidth := fmt.Sprintf("%d", int(*params.W))
			if _, err = os.Stat(dir + "/" + sWidth); os.IsNotExist(err) {
				err = os.MkdirAll(dir+"/"+sWidth, os.ModePerm)
				if err != nil {
					a.internalError(w, err)
					return
				}
			}
			// check files
			files, err := os.ReadDir(dir + "/" + sWidth)
			if err != nil {
				a.service.Logger.Errorln("check files error", dir+"/"+sWidth)
				a.internalError(w, err)
			}

			// if file exists
			if len(files) > 0 {
				for _, f := range files {
					if f.Name() == name+"."+ext {
						a.service.Logger.Infoln("file exists", dir+"/"+sWidth, name+"."+ext)
						a.serveFile(w, dir+"/"+sWidth, name+"."+ext)
						return
					}
				}
			}

			a.service.Logger.Infoln("file not exists, create =>", dir+"/"+sWidth, name+"."+ext)

			// file not exists
			switch ext {
			case "webp":
				err = a.decodeWEBP(src, dir, name+"."+ext, width, height)
				if err != nil {
					a.internalError(w, err)
					return
				}
				a.service.Logger.Infoln("create =>", dir+"/"+sWidth, name+"."+ext)
			case "png":
				err = a.decodePNG(src, dir, name+"."+ext, width, height)
				if err != nil {
					a.internalError(w, err)
					return
				}
				a.service.Logger.Infoln("create =>", dir+"/"+sWidth, name+"."+ext)
			case "jpg":
				err = a.decodeJPG(src, dir, name+"."+ext, width, height)
				if err != nil {
					a.internalError(w, err)
					return
				}
				a.service.Logger.Infoln("create =>", dir+"/"+sWidth, name+"."+ext)
			case "pdf":
				a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
				return
			case "svg":
				a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
				return
			default:
				a.service.Logger.Errorln(name, ext)
				return
			}

			// check file created
			files, err = os.ReadDir(dir + "/" + sWidth)
			if err != nil {
				a.internalError(w, err)
				return
			}
			// if file exists
			if len(files) > 0 {
				for _, f := range files {
					if f.Name() == name+"."+ext {
						a.serveFile(w, dir+"/"+sWidth, name+"."+ext)
						return
					}
				}
			}
		}

		fmt.Printf("ext: %s, width: %d, height: %d\n", ext, width, height)

		switch ext {
		case "webp":
			err = a.decodeBaseWEBP(src, dir, name+"."+ext)
			if err != nil {
				a.internalError(w, err)
			}
		case "png":
			err = a.decodeBasePNG(src, dir, name+"."+ext)
			if err != nil {
				a.internalError(w, err)
			}
		case "jpg":
			err = a.decodeBaseJPG(src, dir, name+"."+ext)
			if err != nil {
				a.internalError(w, err)
			}
		case "pdf":
			a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
			return
		case "svg":
			a.serveFile(w, uploaderApi.UPLOAD_DIR, filename)
			return
		default:
			a.service.Logger.Errorln(name, ext)
			return
		}
		a.serveFile(w, dir, name+"."+ext)
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
	fmt.Printf("width: %d, height: %d\n", width, height)

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
		a.service.Logger.Errorln(err)
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
	buf, err := os.ReadFile(fmt.Sprintf("%s/%s", path, filename))

	if err != nil {
		a.service.Logger.Errorln(err)
	}

	// Decode the image
	img, _, err := image.DecodeConfig(bytes.NewReader(buf))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Extract width and height
	fmt.Printf("Original Width: %d, Height: %d\n", img.Width, img.Height)

	ext := strings.Replace(filepath.Ext(filename), ".", "", 1)

	a.service.Logger.Infoln("serve", filename, ext)

	contentType := ""

	if ext == "svg" {
		contentType = fmt.Sprintf("image/%s+xml", ext)
	}

	if ext == "jpeg" {
		contentType = fmt.Sprintf("image/%s", "jpg")
	}

	if ext == "jpg" {
		contentType = fmt.Sprintf("image/%s", "jpg")
	}

	if ext == "png" {
		contentType = fmt.Sprintf("image/%s", "png")
	}

	if ext == "pdf" {
		contentType = "application/pdf"
	}

	if ext == "webp" {
		contentType = fmt.Sprintf("image/%s", ext)
	}

	if ext == "m4v" {
		contentType = fmt.Sprintf("video/%s", ext)
	}

	if ext == "mp4" {
		contentType = fmt.Sprintf("video/%s", ext)
	}

	if ext == "mov" {
		contentType = fmt.Sprintf("video/%s", ext)
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

// /
func (a *API) serveOriginal(w http.ResponseWriter, dir string, name string, ext string) {
	files, err := os.ReadDir(dir)
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

func (a *API) getSize(src image.Image, pw *float64, pdpr *float64) (int, int) {
	bounds := src.Bounds()

	sourceWidth := bounds.Dx()
	sourceHeight := bounds.Dy()

	fmt.Printf("SourceWidth: %d, SourceHeight: %d\n",
		sourceWidth, sourceHeight,
	)

	width := sourceWidth
	height := sourceHeight

	if pw != nil {

		coef := float64(sourceHeight) / float64(sourceWidth)

		width = int(*pw)
		height = int(coef * *pw)
	}

	dpr := int(1)

	if pdpr != nil {
		dpr = int(*pdpr)
		width = width * dpr
		height = height * dpr
	}

	fmt.Printf("Width: %d, Height: %d\n",
		width, height,
	)

	return width, height
}

func (a *API) getSource(rw http.ResponseWriter, dir string, filename string, ext string) image.Image {
	input, err := os.Open(dir + "/" + filename)
	if err != nil {
		a.internalError(rw, err)
		return nil
	}

	defer func(input *os.File) {
		err = input.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(input)

	sourceExt := a.extractExt(filename)

	fmt.Printf("sourceExt %s, ext %s\n", sourceExt, ext)

	switch sourceExt {
	case "svg":
		src, err := a.decodeBaseSVG(input)
		if err != nil {
			a.internalError(rw, err)
			return nil
		}
		return src
	default:
		// decode source file
		src, _, err := image.Decode(input)
		if err != nil {
			a.internalError(rw, err)
			return nil
		}
		return src
	}
}

func (a *API) decodeBaseSVG(r io.Reader) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(r)
	if err != nil {
		return nil, err
	}

	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	icon.SetTarget(0, 0, float64(w), float64(h))

	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	return rgba, err
}
