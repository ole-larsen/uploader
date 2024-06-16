package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/ole-larsen/uploader/models"
	"github.com/ole-larsen/uploader/service/db"
)

const PublicDir = "/api/v1/files/"

type File struct {
	ID        int64       `db:"id"`
	Name      string      `db:"name"`
	Alt       string      `db:"alt"`
	Caption   string      `db:"caption"`
	Width     int64       `db:"width"`
	Height    int64       `db:"height"`
	Provider  *string     `db:"provider"`
	Hash      string      `db:"hash"`
	Ext       string      `db:"ext"`
	Size      float64     `db:"size"`
	Url       string      `db:"url"`
	Formats   interface{} `db:"formats"`
	Metadata  interface{} `db:"metadata"`
	Mime      string      `db:"mime"`
	Thumb     string      `db:"preview_url"`
	CreatedBy int64       `db:"created_by"`
	UpdatedBy int64       `db:"updated_by"`
	Created   strfmt.Date `db:"created"`
	Updated   strfmt.Date `db:"updated"`
	Deleted   strfmt.Date `db:"deleted"`
}

type FileRepository interface {
	Create(fileMap map[string]interface{}) error
	Update(fileMap map[string]interface{}) ([]*models.File, error)
	GetFiles() ([]*models.File, error)
	GetFileByName(name string) (*models.File, error)
	GetFileByID(id int64) (*models.File, error)
	GetPublicFilesByProvider(_provider string) ([]*models.PublicFile, error)
	GetPublicFileByName(name string) (*models.PublicFile, error)
	GetPublicFileByID(id int64) (*models.PublicFile, error)
}

type FileRepo struct {
	tbl       string
	db        *sqlx.DB
	ctx       context.Context
	storage   string
	PublicDir string
}

func NewFileRepo(ctx context.Context, store db.Storer) *FileRepo {
	return &FileRepo{
		tbl:       "files",
		db:        store.InnerDB(),
		ctx:       ctx,
		PublicDir: PublicDir,
	}
}

func (f FileRepo) Create(fileMap map[string]interface{}) error {
	if f.db == nil {
		return errDbNotInitialized
	}

	_, err := f.db.NamedExec(`
		INSERT INTO files (name, alt, caption, hash, mime, ext, size, width, height, provider, url)
		VALUES (:name, :alt, :caption, :hash, :mime, :ext, :size, :width, :height, :provider, :url)
		ON CONFLICT DO NOTHING`, fileMap)
	return err
}

func (f FileRepo) Update(fileMap map[string]interface{}) ([]*models.File, error) {
	if f.db == nil {
		return nil, errDbNotInitialized
	}
	fileMap["url"] = fmt.Sprintf("%s%s%s", PublicDir, fileMap["name"], fileMap["ext"])
	_, err := f.db.NamedExec(`UPDATE files SET
                name=:name,
                alt=:alt,
                hash=:hash,
                caption=:caption,
                ext=:ext,
                mime=:mime,
                size=:size,
                width=:width,
                height=:height,
                url=:url,
                provider=:provider WHERE id =:id`, fileMap)
	if err != nil {
		return nil, err
	}
	return f.GetFiles()
}

func (f FileRepo) GetFiles() ([]*models.File, error) {
	if f.db == nil {
		return nil, errDbNotInitialized
	}
	var (
		multierr multierror.Error
		files    []*models.File
	)

	rows, err := f.db.Queryx(
		"SELECT id, name, alt, caption, hash, mime, ext, size, width, height, provider, url, created, updated, deleted from files;")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var file File
		err = rows.Scan(&file.ID, &file.Name, &file.Alt, &file.Caption, &file.Hash, &file.Mime, &file.Ext, &file.Size,
			&file.Width, &file.Height, &file.Provider, &file.Url, &file.Created, &file.Updated, &file.Deleted)
		if err != nil {
			return nil, err
		}

		provider := ""
		if file.Provider != nil {
			provider = *file.Provider
		}
		var size int64 = int64(math.Round(file.Size))

		files = append(files, &models.File{
			ID:       file.ID,
			Name:     file.Name + file.Ext,
			Thumb:    file.Url,
			Alt:      file.Alt,
			Caption:  file.Caption,
			Hash:     file.Hash,
			Type:     file.Mime,
			Ext:      file.Ext,
			Size:     size,
			Width:    file.Width,
			Height:   file.Height,
			Provider: provider,
			Created:  file.Created,
			Updated:  file.Updated,
			Deleted:  file.Deleted,
		})
	}
	defer rows.Close()

	return files, multierr.ErrorOrNil()
}

func (f FileRepo) GetFileByName(name string) (*models.File, error) {
	if f.db == nil {
		return nil, errDbNotInitialized
	}
	var file File
	sqlStatement := `SELECT id, name, alt, caption, hash, mime, ext, size, width, height, provider, url, created, updated, deleted from files WHERE name=$1;`
	row := f.db.QueryRow(sqlStatement, name)
	err := row.Scan(&file.ID, &file.Name, &file.Alt, &file.Caption, &file.Hash, &file.Mime, &file.Ext, &file.Size,
		&file.Width, &file.Height, &file.Provider, &file.Url, &file.Created, &file.Updated, &file.Deleted)
	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("file not found")
	case nil:
		provider := ""
		if file.Provider != nil {
			provider = *file.Provider
		}
		var size int64 = int64(math.Round(file.Size))

		return &models.File{
			ID:       file.ID,
			Name:     file.Name + file.Ext,
			Thumb:    file.Url,
			Alt:      file.Alt,
			Caption:  file.Caption,
			Hash:     file.Hash,
			Type:     file.Mime,
			Ext:      file.Ext,
			Size:     size,
			Width:    file.Width,
			Height:   file.Height,
			Provider: provider,
			Created:  file.Created,
			Updated:  file.Updated,
			Deleted:  file.Deleted,
		}, err
	default:
		return nil, err
	}
}

func (f FileRepo) GetFileByID(id int64) (*models.File, error) {
	if f.db == nil {
		return nil, errDbNotInitialized
	}
	var file File
	sqlStatement := `SELECT id, name, alt, caption, hash, mime, ext, size, width, height, provider, url, created, updated, deleted from files WHERE id=$1;`
	row := f.db.QueryRow(sqlStatement, id)
	err := row.Scan(&file.ID, &file.Name, &file.Alt, &file.Caption, &file.Hash, &file.Mime, &file.Ext, &file.Size,
		&file.Width, &file.Height, &file.Provider, &file.Url, &file.Created, &file.Updated, &file.Deleted)
	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("file not found")
	case nil:
		provider := ""
		if file.Provider != nil {
			provider = *file.Provider
		}
		var size int64 = int64(math.Round(file.Size))

		return &models.File{
			ID:       file.ID,
			Name:     file.Name,
			Thumb:    file.Url,
			Alt:      file.Alt,
			Caption:  file.Caption,
			Hash:     file.Hash,
			Type:     file.Mime,
			Ext:      file.Ext,
			Size:     size,
			Width:    file.Width,
			Height:   file.Height,
			Provider: provider,
			Created:  file.Created,
			Updated:  file.Updated,
			Deleted:  file.Deleted,
		}, err
	default:
		return nil, err
	}
}

func (f FileRepo) GetPublicFilesByProvider(_provider string) ([]*models.PublicFile, error) {
	if f.db == nil {
		return nil, errDbNotInitialized
	}
	var (
		multierr multierror.Error
		files    []*models.PublicFile
	)

	rows, err := f.db.Queryx(
		`SELECT id, name, alt, caption, hash, mime, ext, size, width, height, provider, url,
   created, updated, deleted from files where provider=$1;`, _provider)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var file File
		err = rows.Scan(&file.ID, &file.Name, &file.Alt, &file.Caption, &file.Hash, &file.Mime, &file.Ext, &file.Size,
			&file.Width, &file.Height, &file.Provider, &file.Url, &file.Created, &file.Updated, &file.Deleted)
		if err != nil {
			return nil, err
		}

		provider := ""
		if file.Provider != nil {
			provider = *file.Provider
		}
		var size int64 = int64(math.Round(file.Size))

		files = append(files, &models.PublicFile{
			ID: file.ID,
			Attributes: &models.PublicFileAttributes{
				Alt:      file.Alt,
				Caption:  file.Caption,
				Ext:      file.Ext,
				Hash:     file.Hash,
				Height:   file.Height,
				Mime:     file.Mime,
				Name:     file.Name,
				Provider: provider,
				Size:     size,
				URL:      file.Url,
				Width:    file.Width,
				Created:  file.Created,
				Updated:  file.Updated,
				Deleted:  file.Deleted,
			},
		})
	}
	defer rows.Close()

	return files, multierr.ErrorOrNil()
}

func (f FileRepo) GetPublicFileByName(name string) (*models.PublicFile, error) {
	file, err := f.GetFileByName(name)
	if err != nil {
		return nil, err
	}
	return &models.PublicFile{
		ID: file.ID,
		Attributes: &models.PublicFileAttributes{
			Alt:      file.Alt,
			Caption:  file.Caption,
			Ext:      file.Ext,
			Hash:     file.Hash,
			Height:   file.Height,
			Mime:     file.Type,
			Name:     file.Name,
			Provider: file.Provider,
			Size:     file.Size,
			URL:      file.Thumb,
			Width:    file.Width,
			Created:  file.Created,
			Updated:  file.Updated,
			Deleted:  file.Deleted,
		},
	}, nil
}

func (f FileRepo) GetPublicFileByID(id int64) (*models.PublicFile, error) {
	file, err := f.GetFileByID(id)
	if err != nil {
		return nil, err
	}
	return &models.PublicFile{
		ID: file.ID,
		Attributes: &models.PublicFileAttributes{
			Alt:      file.Alt,
			Caption:  file.Caption,
			Ext:      file.Ext,
			Hash:     file.Hash,
			Height:   file.Height,
			Mime:     file.Type,
			Name:     file.Name,
			Provider: file.Provider,
			Size:     file.Size,
			URL:      file.Thumb,
			Width:    file.Width,
			Created:  file.Created,
			Updated:  file.Updated,
			Deleted:  file.Deleted,
		},
	}, nil
}
