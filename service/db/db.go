package db

import (
	"errors"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	postgresProto        = "postgres"
	mySQLProto           = "mysql"
	MaxMySQLPlaceholders = 65535
)

type Storer interface {
	Type() string
	InnerDB() *sqlx.DB
	Close() error
}

type database struct {
	*sqlx.DB
	statsPollerCloseChan chan struct{}
	dbType               string
}

func SetupRDBMS(dsn string) (_ Storer, err error) {
	var (
		db *sqlx.DB
	)

	splitted := strings.Split(dsn, "://")
	if len(splitted) < 2 {
		return &database{}, errors.New("can't parse DSN")
	}
	proto := splitted[0]
	source := splitted[1]
	switch proto {
	case mySQLProto:
		db, err = sqlx.Connect(mySQLProto, source)
	case postgresProto:
		db, err = sqlx.Connect(postgresProto, dsn)
		if err != nil {
			err = errors.New(strings.Replace(err.Error(), dsn, "DSN", -1))
		}
	default:
		return &database{}, errors.New("unknown database protocol")
	}
	if err != nil {
		return &database{}, err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(25)
	db.SetConnMaxIdleTime(time.Hour)
	db.SetConnMaxLifetime(time.Hour)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	statsPollerCloseChan := make(chan struct{})
	
	return &database{
		DB:                   db,
		dbType:               proto,
		statsPollerCloseChan: statsPollerCloseChan,
	}, nil

}

func (d *database) Type() string {
	if d == nil {
		return ""
	}
	return d.dbType
}

func (d *database) InnerDB() *sqlx.DB {
	if d == nil {
		return nil
	}
	return d.DB
}

func (d *database) Close() error {
	if d == nil {
		return nil
	}
	if d.DB == nil {
		return nil
	}
	close(d.statsPollerCloseChan)
	return d.DB.Close()
}
