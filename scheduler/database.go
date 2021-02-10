package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// DB holds redis conn
type DB struct {
	Conn redis.Conn
}

// ConnectDB constructor function for Redis db
func ConnectDB(url string) (db DB, err error) {
	fmt.Println(url)
	conn, err := redis.DialURL(url)
	if err != nil {
		fmt.Println("error connecting to db")
	}

	db = DB{
		Conn: conn,
	}

	return db, err
}
