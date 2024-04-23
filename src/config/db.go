package config

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB
var e error

func DatabaseInit() {
	defaultHost := "localhost"
	defaultPort := 5432
	user := "root"
	password := "password"
	dbName := "go_rest_api"

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}
	portStr := os.Getenv("DB_PORT")
	var port int
	if portStr == "" {
		port = defaultPort
	} else {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			panic(err)
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e != nil {
		panic(e)
	}
}

func DB() *gorm.DB {
	return database
}
