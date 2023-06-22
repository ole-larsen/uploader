package uploaderApi

import "github.com/ole-larsen/uploader/service"

const UPLOAD_DIR = "./uploads"
const UPLOAD_DIR_RAW = "./uploads/raw"

type API struct {
	service *service.Service
}

func NewUploaderAPI(s *service.Service) *API {
	return &API{service: s}
}
