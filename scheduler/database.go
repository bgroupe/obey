package main

import (
	"github.com/gomodule/redigo/redis"

	log "github.com/sirupsen/logrus"
)

// DB holds redis conn
type DB struct {
	Conn redis.Conn
}

// ConnectDB constructor function for Redis db
func ConnectDB(url string) (db DB, err error) {
	log.Info("connected to redis")
	conn, err := redis.DialURL(url)
	if err != nil {
		log.Error("error connecting to db")
	}

	db = DB{
		Conn: conn,
	}

	return db, err
}
