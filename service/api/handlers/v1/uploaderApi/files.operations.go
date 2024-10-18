package uploaderApi

import (
	"crypto/sha256"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/buckket/go-blurhash"
	"github.com/chai2010/webp"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/restapi/operations/uploader"
	"github.com/ole-larsen/uploader/service/db/repository"
	"github.com/ole-larsen/uploader/service/settings"
)

func (a *API) postFiles(params uploader.PostUploaderFilesParams, _ *models.Principal) (*models.PublicFile, error) {
	file, fileHeader, err := params.HTTPRequest.FormFile("file")

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(file)

	if err != nil {
		return nil, err
	}

	attributes := make(map[string]interface{})

	ext := filepath.Ext(fileHeader.Filename)

	if params.HTTPRequest.Form.Get("ext") != "" && params.HTTPRequest.Form.Get("ext") != "undefined" {
		ext = params.HTTPRequest.Form.Get("ext")
	}

	filenameWithExtension := params.HTTPRequest.Form.Get("name")

	filename := strings.TrimSuffix(filenameWithExtension, ext)

	// modify filename
	if settings.Settings.UseDB {
		tag := time.Now().Format(time.RFC850)
		filename = filename + tag
	}
	// Here we start with a new hash.
	h := sha256.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(filename))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	hash := fmt.Sprintf("%x", bs)

	attributes["hash"] = params.HTTPRequest.Form.Get("hash")

	if params.HTTPRequest.Form.Get("hash") == "" {
		attributes["hash"] = hash
	}

	attributes["name"] = filename

	attributes["alt"] = params.HTTPRequest.Form.Get("alt")
	attributes["caption"] = params.HTTPRequest.Form.Get("caption")
	attributes["mime"] = params.HTTPRequest.Form.Get("type")
	attributes["size"] = params.HTTPRequest.Form.Get("size")
	attributes["width"] = params.HTTPRequest.Form.Get("width")
	attributes["height"] = params.HTTPRequest.Form.Get("height")
	attributes["ext"] = ext
	attributes["provider"] = params.HTTPRequest.Form.Get("provider")
	attributes["blur"] = params.HTTPRequest.Form.Get("blur")

	attributes["url"] = fmt.Sprintf("%s%s%s", repository.PublicDir, attributes["name"], attributes["ext"])

	if settings.Settings.UseHash {
		attributes["url"] = fmt.Sprintf("%s%s%s", repository.PublicDir, hash, attributes["ext"])
	}

	if settings.Settings.UseDB {
		exists, err := a.service.Files.GetFileByName(filename)

		if err != nil && err.Error() != "file not found" {
			return nil, err
		}

		if exists != nil {
			return nil, fmt.Errorf("file exists")
		}
	}

	fmt.Println("--------------------------------")
	fmt.Println(attributes)
	fmt.Println("--------------------------------")

	if _, err = os.Stat(UPLOAD_DIR); os.IsNotExist(err) {
		err = os.MkdirAll(UPLOAD_DIR, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		files, err := os.ReadDir(UPLOAD_DIR)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			if f.Name() == fmt.Sprintf("%s%s", filename, ext) {
				return nil, fmt.Errorf("file exists")
			}
		}
	}
	name, _ := attributes["name"].(string)

	if settings.Settings.UseHash {
		_, err = a.createFile(file, UPLOAD_DIR, hash, ext)
		if err != nil {
			fmt.Println("OK store to disk", err)
			return nil, err
		}
	} else {
		_, err = a.createFile(file, UPLOAD_DIR, name, ext)
		fmt.Println("OK store to disk", err)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("EXTENSION FROM FILE", ext)
	switch ext {
	case ".webp":
		loadedImage, err := webp.Decode(file)
		if err != nil {
			fmt.Println("error from webp decode", err)
			return nil, err
		}
		str, _ := blurhash.Encode(4, 3, loadedImage)
		if err != nil {
			fmt.Println("error from blurhash encode", err)
		}
		fmt.Printf("Hash: %s\n", str)
		attributes["blur"] = str
	case ".png":
		loadedImage, err := png.Decode(file)
		if err != nil {
			fmt.Println("error from png decode", err)
			return nil, err
		}
		str, _ := blurhash.Encode(4, 3, loadedImage)
		if err != nil {
			fmt.Println("error from blurhash encode", err)
		}
		fmt.Printf("Hash: %s\n", str)
		attributes["blur"] = str
	case ".jpg", ".jpeg":
		loadedImage, err := jpeg.Decode(file)
		if err != nil {
			fmt.Println("error from jpg decode", err)
			return nil, err
		}
		str, _ := blurhash.Encode(4, 3, loadedImage)
		if err != nil {
			fmt.Println("error from blurhash encode", err)
		}
		fmt.Printf("Hash: %s\n", str)
		attributes["blur"] = str
	case ".pdf", ".svg":
	default:
		a.service.Logger.Errorln(name, ext)
	}

	if settings.Settings.UseDB {
		if err = a.service.Files.Create(attributes); err != nil {
			fmt.Println("OK store to database", err)
			return nil, err
		}
	}

	if settings.Settings.UseDB {
		return a.service.Files.GetPublicFileByName(filename)
	} else {
		return &models.PublicFile{
			ID: 0,
			Attributes: &models.PublicFileAttributes{
				Name: name,
				Hash: hash,
			},
		}, nil
	}
}

func (a *API) createFile(file multipart.File, directory string, name string, ext string) (*os.File, error) {
	a.service.Logger.Infoln(fmt.Sprintf("create %s/%s%s", directory, name, ext))
	dst, err := os.Create(fmt.Sprintf("%s/%s%s", directory, name, ext))
	if err != nil {
		return nil, err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)

	return dst, err
}

func (a *API) putFiles(params uploader.PutUploaderFilesParams, principal *models.Principal) ([]*models.File, error) {
	file, fileHeader, err := params.HTTPRequest.FormFile("file")

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(file)

	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(fileHeader.Filename)
	if params.HTTPRequest.Form.Get("ext") != "" {
		ext = params.HTTPRequest.Form.Get("ext")
	}

	filename := params.HTTPRequest.Form.Get("name")

	id, err := strconv.ParseInt(params.HTTPRequest.Form.Get("id"), 10, 64)

	if err != nil {
		return nil, err
	}

	attributes := make(map[string]interface{})
	attributes["id"] = id
	attributes["hash"] = params.HTTPRequest.Form.Get("hash")
	attributes["name"] = filename
	attributes["ext"] = ext
	attributes["alt"] = params.HTTPRequest.Form.Get("alt")
	attributes["caption"] = params.HTTPRequest.Form.Get("caption")
	attributes["mime"] = params.HTTPRequest.Form.Get("type")
	attributes["size"] = params.HTTPRequest.Form.Get("size")
	attributes["width"] = params.HTTPRequest.Form.Get("width")
	attributes["height"] = params.HTTPRequest.Form.Get("height")
	attributes["created_by_id"] = *principal
	attributes["updated_by_id"] = *principal
	attributes["provider"] = params.HTTPRequest.Form.Get("provider")
	attributes["blur"] = params.HTTPRequest.Form.Get("blur")
	exist, err := a.service.Files.GetFileByID(id)

	if err != nil {
		return nil, err
	}

	path := strings.Split(exist.Thumb, "/")
	encodedFilename := path[len(path)-1]

	existFilename, err := url.QueryUnescape(encodedFilename)
	if err != nil {
		return nil, err
	}

	err = os.Remove(fmt.Sprintf("%s/%s", UPLOAD_DIR, existFilename))
	if err != nil {
		return nil, err
	}

	err = os.Remove(fmt.Sprintf("%s/%s%s", UPLOAD_DIR, filename, ext))
	if err != nil {
		return nil, err
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("%s/%s%s", UPLOAD_DIR, filename, ext))

	if err != nil {
		return nil, err
	}

	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			a.service.Logger.Errorln(err)
		}
	}(dst)

	// Copy the uploaded file to the filesystem at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}
	fmt.Println("PUT file:", attributes)
	return a.service.Files.Update(attributes)
}

func (a *API) getFilesID(params uploader.GetUploaderFilesIDParams, _ *models.Principal) (*models.File, error) {
	return a.service.Files.GetFileByID(params.ID)
}

func (a *API) getFiles(params uploader.GetUploaderFilesParams, _ *models.Principal) ([]*models.File, error) {
	if params.Name != nil {
		name := *params.Name
		ext := filepath.Ext(name)

		file, err := a.service.Files.GetFileByName(strings.TrimSuffix(name, ext))
		if err != nil {
			return nil, err
		}
		return append(make([]*models.File, 0), file), nil
	}
	return a.service.Files.GetFiles()
}
