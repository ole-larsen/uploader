package service

import (
	"context"

	"github.com/ole-larsen/uploader/service/db"
	"github.com/ole-larsen/uploader/service/db/repository"
	"github.com/ole-larsen/uploader/service/log"
	"github.com/ole-larsen/uploader/service/settings"
)

type channels struct {
	cancel context.CancelFunc
	done   chan struct{}
}

type Repositories struct {
	Files repository.FileRepository
}

type Service struct {
	channels
	Repositories
	db     db.Storer
	Logger log.Logger
	Ctx    context.Context
}

func NewService() *Service {
	logger := log.NewLogger()

	ctx, cancel := context.WithCancel(context.Background())

	service := &Service{
		channels: channels{
			cancel: cancel,
			done:   make(chan struct{}),
		},
		Logger: logger,
		Ctx:    ctx,
	}

	if settings.Settings.UseDB == true {
		psqlDB, err := db.SetupRDBMS(settings.Settings.PGSQL)

		if err != nil {
			logger.Fatalln("postgres init:", err)
		}

		files := repository.NewFileRepo(ctx, psqlDB)
		service.Repositories = Repositories{
			Files: files,
		}
		service.db = psqlDB

	}

	return service
}
