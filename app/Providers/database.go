package services

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jinzhu/gorm"

	//Framework
	c "../../config"
)

var dbConnections = make(map[string]*sql.DB)
var dbMutex sync.RWMutex

func SqlDB(ctx context.Context, connectionTag ...string) (*sql.DB, error) {

	var connectionString string
	if len(connectionTag) == 0 {
		connectionString = c.DEFAULT_CONNECTION_TAG //Use default database connection
	} else {
		connectionString = connectionTag[0]
	}

	var db *sql.DB

	// Create a retryable operation
	operation := func() error {

		//Check if database connection already in dbConnections Map
		dbMutex.RLock()
		connection := dbConnections[connectionString]
		dbMutex.RUnlock()

		justCreated := false // Used to indicate if connection was registered for first time

		if connection == nil {
			//Connection does not exist - make a new one
			var err error
			connection, err = sql.Open(connectionOpenString(connectionString))
			if err != nil {
				return err //Retry attempt
			}
			d := c.Connections(connectionString)
			connection.SetConnMaxLifetime(d.SetConnMaxLifetime)
			connection.SetMaxIdleConns(d.SetMaxIdleConns)
			connection.SetMaxOpenConns(d.SetMaxOpenConns)
			justCreated = true
		}

		// Call PingContext here and test for driver.ErrBadConn. If driver.ErrBadConn, don't try again. Otherwise try again.
		err := connection.PingContext(ctx)
		if err != nil {

			dbMutex.Lock()
			connection.Close()
			if !justCreated {
				//Remove from dbConnections
				delete(dbConnections, connectionString)
			}
			dbMutex.Unlock()
			if err == driver.ErrBadConn {
				return backoff.Permanent(err) //Don't retry
			} else {
				return err //Retry
			}
		}

		if justCreated {
			//Store into dbConnections Map so we can reuse again
			dbMutex.Lock()
			dbConnections[connectionString] = connection
			dbMutex.Unlock()
		}

		db = connection
		return nil //All good
	}

	backoffAlgorithm := backoff.NewExponentialBackOff()
	backoffAlgorithm.MaxElapsedTime = time.Duration(10000) * time.Millisecond
	err := backoff.Retry(operation, backoffAlgorithm)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// DebugLogger is a custom logger for Gorm because historically it had issues logging output whilst in local development.
// This `hack` seemed to fix the problem.
type DebugLogger struct {
}

func (w DebugLogger) Print(v ...interface{}) {
	log.Print(v)
}

// Return a gorm object with the relevant underlying *sql.DB object
// NOTE: Gorm internally calls Ping(). TODO: Use a fork of Gorm that removes unnecessary Ping() call
// since SqlDB() already calls Ping()
func Gorm(ctx context.Context, connectionTag ...string) (*gorm.DB, error) {

	db, err := SqlDB(ctx, connectionTag...)
	if err != nil {
		return nil, err
	}

	var connectionString string
	if len(connectionTag) == 0 {
		connectionString = c.DEFAULT_CONNECTION_TAG //Use default database connection
	} else {
		connectionString = connectionTag[0]
	}
	dbName, _ := connectionOpenString(connectionString)

	gormDB, err := gorm.Open(dbName, db)
	if err != nil {
		// NOTE: It is possible that SqlDB() provided us with a healthy connection.
		// Internally Gorm calls Ping() too, so it is possible that the connection became stale in the meantime.
		// Let's remove it from our store of connections since it's stale
		if err == driver.ErrBadConn {
			dbMutex.Lock()
			connection := dbConnections[connectionString]
			if connection != nil {
				connection.Close()
			}
			delete(dbConnections, connectionString) //Remove from dbConnections
			dbMutex.Unlock()
		}

		return nil, err
	}

	// Set for Local Development
	// gormDB.LogMode(true)
	// gormDB.SetLogger(DebugLogger{})

	return gormDB, nil
}

func connectionOpenString(connectionTag ...string) (string, string) {

	var db *c.Database

	if len(connectionTag) == 0 {
		db = c.Connections(c.DEFAULT_CONNECTION_TAG)
		if db == nil {
			panic("No valid connectionTag specified - DEFAULT_CONNECTION_TAG setting in config/database.go is invalid")
		}
	} else if len(connectionTag) == 1 {
		db = c.Connections(connectionTag[0])
		if db == nil {
			panic("No valid connectionTag specified")
		}
	} else {
		panic("Too many arguments - connectionOpenString requires 0 or 1 argument")
	}

	if db.Settings == "" {
		if db.Password == "" {
			return db.Driver, db.User + "@" + db.Protocol + "(" + db.Host + ":" + db.Port + ")/" + db.Name
		} else {
			return db.Driver, db.User + ":" + db.Password + "@" + db.Protocol + "(" + db.Host + ":" + db.Port + ")/" + db.Name
		}
	} else {
		if db.Password == "" {
			return db.Driver, db.User + "@" + db.Protocol + "(" + db.Host + ":" + db.Port + ")/" + db.Name + "?" + db.Settings
		} else {
			return db.Driver, db.User + ":" + db.Password + "@" + db.Protocol + "(" + db.Host + ":" + db.Port + ")/" + db.Name + "?" + db.Settings
		}
	}
}
