package config

import ()

var (
	MYSQL_HOST     = "localhost"
	MYSQL_PORT     = "3306"
	MYSQL_PROTOCOL = "tcp"
	MYSQL_USER     = "root"
	MYSQL_PASSWORD = ""
)

/*
|--------------------------------------------------------------------------
| Default Database Connection Name
|--------------------------------------------------------------------------
|
| Here you may specify which of the database connections below you wish
| to use as your default connection for all database work. Of course
| you may use many connections at once using the Database library.
|
*/

const DEFAULT_CONNECTION_TAG = "default"

/*
|--------------------------------------------------------------------------
| Database Connections
|--------------------------------------------------------------------------
|
| Here are each of the database connections setup for your application.
|
|
*/

type database struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Protocol string
	Settings string
}

func Connections(connectionTag string) *database {

	connections := make(map[string]*database)

	connections["default"] = &database{
		Driver:   "mysql",
		Host:     MYSQL_HOST,
		Port:     MYSQL_PORT,
		Name:     "",
		User:     MYSQL_USER,
		Password: MYSQL_PASSWORD,
		Protocol: MYSQL_PROTOCOL,
		Settings: "parseTime=true", //"?parseTime=true"
	}

	return connections[connectionTag]
}
