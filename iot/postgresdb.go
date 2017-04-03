package iot

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-martini/martini"
	// Using blank import for Postgres Driver
	_ "github.com/lib/pq"
)

// PostgresSession is the struct to keep Sarama
type PostgresSession struct {
	*sql.DB
	postgresConnStr string
}

// NewPostgresSession connects to the Postgres
func NewPostgresSession(connStr string) *PostgresSession {

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	return &PostgresSession{db, connStr}
}

// RegisterDevice is the utility function to register a device from service
func (postgresSession *PostgresSession) RegisterDevice(device Device) error {
	// first check whether the device is already in the Database
	// if so do not do anything, if not add it to the db

	err := deviceAlreadyRegistered(postgresSession.DB, device)

	if err != nil {
		return err
	}

	return insertDevice(postgresSession.DB, device)
}

func deviceAlreadyRegistered(postgres *sql.DB, device Device) error {
	// TODO: gseng - SQL Injection here :-)
	queryDeviceStr := `SELECT * FROM "Device" WHERE "UniqueDeviceId" = $1`
	deviceResult := Device{}

	rows, err := postgres.Query(queryDeviceStr, device.UniqueDeviceID)

	if err != nil {
		return err
	}

	values, err := rows.Columns()

	if err != nil {
		return err
	}

	if len(values) == 1 {
		fmt.Printf("Device with id %s already registered\n", device.UniqueDeviceID)
		return nil
	} else if len(values) > 0 {
		return fmt.Errorf("More than one device found with the same id %s", device.UniqueDeviceID)
	}

	fmt.Printf("Returned device id is %v\n", deviceResult.UniqueDeviceID)

	return nil
}

func insertDevice(postgres *sql.DB, device Device) error {
	sqlStatement := `  
		INSERT INTO "Device" ("UniqueDeviceId", "DeviceType", "CreatedOn", "CreatedBy", "UpdatedOn")  
		VALUES ($1, $2, $3, $4, $5)
		RETURNING "Id"`

	var datetime = time.Now()
	datetime.Format(time.RFC3339)

	createdBy := "00000000-0000-0000-0000-000000000001"

	id := 0
	err := postgres.QueryRow(sqlStatement, device.UniqueDeviceID, device.DeviceType, datetime, createdBy, datetime).Scan(&id)
	if err != nil {
		return err
	}

	fmt.Printf("New record ID is: %d", id)

	return nil
}

// NewPostgresSession adds Postgres to the Martini pipeline
func (postgresSession *PostgresSession) NewPostgresSession() martini.Handler {
	return func(context martini.Context) {
		context.Map(postgresSession)
		context.Next()
	}
}
