package config

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	MYSQL_HOST     = "localhost"
	MYSQL_PORT     = "3306"
	MYSQL_PROTOCOL = "tcp"
	MYSQL_USER     = "root"
	MYSQL_PASSWORD = ""
)

// DEFAULT_CONNECTION_TAG sets which connection below is the default connection for use by
const DEFAULT_CONNECTION_TAG = "default"

type Database struct {
	Driver             string
	Host               string
	Port               string
	Name               string
	User               string
	Password           string
	Protocol           string
	Settings           string
	SetConnMaxLifetime time.Duration
	SetMaxIdleConns    int //Always set this to prevent it from being zero
	SetMaxOpenConns    int
}

func Connections(connectionTag string) *Database {

	connections := make(map[string]*Database)

	connections["default"] = &Database{
		Driver:          "mysql",
		Host:            MYSQL_HOST,
		Port:            MYSQL_PORT,
		Name:            "",
		User:            MYSQL_USER,
		Password:        MYSQL_PASSWORD,
		Protocol:        MYSQL_PROTOCOL,
		Settings:        "parseTime=true",
		SetMaxIdleConns: 10,
		SetMaxOpenConns: 12,
	}

	return connections[connectionTag]
}
