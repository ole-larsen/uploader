package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table files")
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS files (
				id SERIAL PRIMARY KEY,
				name character varying(255) UNIQUE NOT NULL,
				alt character varying(128),
				caption character varying(255),
				width integer,
				height integer,
				formats jsonb,
				hash character varying(128),
				ext character varying(10),
				mime character varying(128),
				size integer,
				url character varying(255),
				preview_url character varying(255),
				provider character varying(128),
				provider_metadata jsonb,
			    created_by_id integer,
				updated_by_id integer,
				created             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				updated             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				deleted             TIMESTAMP WITH TIME ZONE DEFAULT NULL
			);`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping files")
		_, err := db.Exec(`DROP TABLE IF EXISTS files`)
		return err
	})
}
