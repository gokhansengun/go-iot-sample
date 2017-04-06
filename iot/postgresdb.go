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

// FetchAlDevices retrieves all the devices registered so far in the db
func (postgresSession *PostgresSession) FetchAlDevices() ([]Device, error) {
	queryAllDevicesStr := `SELECT "UniqueDeviceId", "DeviceType" FROM "Device"`

	rows, err := postgresSession.DB.Query(queryAllDevicesStr)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	uniqueDeviceID := ""
	deviceType := -1
	devices := []Device{}
	for rows.Next() {
		err = rows.Scan(&uniqueDeviceID, &deviceType)

		device := Device{UniqueDeviceID: uniqueDeviceID, DeviceType: deviceType}

		devices = append(devices, device)

		if err != nil {
			return nil, err
		}
	}

	return devices, nil
}

// RegisterDevice is the utility function to register a device from service
func (postgresSession *PostgresSession) RegisterDevice(device Device) error {
	// first check whether the device is already in the Database
	// if so do not do anything, if not add it to the db

	alreadyRegistered, err := deviceAlreadyRegistered(postgresSession.DB, device.UniqueDeviceID)

	if err != nil {
		return err
	}

	if !alreadyRegistered {
		return insertDevice(postgresSession.DB, device)
	}

	return nil
}

// DoesDeviceExist checks whether device is already in the db or not
func (postgresSession *PostgresSession) DoesDeviceExist(uniqueDeviceID string) bool {
	alreadyExists, _ := deviceAlreadyRegistered(postgresSession.DB, uniqueDeviceID)

	return alreadyExists
}

func deviceAlreadyRegistered(postgres *sql.DB, uniqueDeviceID string) (bool, error) {
	// TODO: gseng - SQL Injection here :-)
	queryDeviceStr := `SELECT COUNT("Id") FROM "Device" WHERE "UniqueDeviceId" = $1`

	rows, err := postgres.Query(queryDeviceStr, uniqueDeviceID)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	count := 0
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, err
		}
	}

	if count == 1 {
		// fmt.Printf("Device with id %s already registered\n", uniqueDeviceID)
		return true, nil
	} else if count > 0 {
		return true, fmt.Errorf("More than one device found with the same id %s", uniqueDeviceID)
	}

	return false, nil
}

func insertDevice(postgres *sql.DB, device Device) error {
	sqlStatement := `  
		INSERT INTO "Device" ("UniqueDeviceId", "DeviceType", "StatusType", "CreatedOn", "CreatedBy", "UpdatedOn", "UpdatedBy")  
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "Id"`

	var datetime = time.Now()
	datetime.Format(time.RFC3339)

	createdBy := "00000000-0000-0000-0000-000000000001"
	updatedBy := "00000000-0000-0000-0000-000000000001"

	id := 0
	err := postgres.QueryRow(sqlStatement, device.UniqueDeviceID, device.DeviceType, 1, datetime, createdBy, datetime, updatedBy).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// NewPostgresHandler adds Postgres to the Martini pipeline
func (postgresSession *PostgresSession) NewPostgresHandler() martini.Handler {
	return func(context martini.Context) {
		context.Map(postgresSession)
		context.Next()
	}
}
