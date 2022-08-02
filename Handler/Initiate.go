package handler

import (
	"database/sql"
	model "finalexam/Model"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

var cfg model.ConfigData
var err = cleanenv.ReadConfig(".config.yml", &cfg)

func dbConn() *sql.DB {
	var err error
	req, err := sql.Open(cfg.Host, "server="+cfg.Server+";database="+cfg.Database+";trusted_connection=yes")
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	return req
}

var db = dbConn()
