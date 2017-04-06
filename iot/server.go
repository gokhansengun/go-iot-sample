package iot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"strconv"

	"time"

	"fmt"

	"github.com/go-martini/martini"
	"github.com/gokhansengun/go-iot-sample/utility"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

// MartiniServer composed of the internal structure and the router
type MartiniServer struct {
	*martini.Martini
	martini.Router
}

func initMartiniServer() *MartiniServer {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &MartiniServer{m, r}
}

// NewServer create the server and set up middleware.
func NewServer(session *DatabaseSession, kafkaSession *KafkaSession, postgresSession *PostgresSession, redisSession *RedisSession) *MartiniServer {
	m := initMartiniServer()
	m.Martini.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Martini.Use(session.NewMongoDbHandler())
	m.Martini.Use(kafkaSession.NewKafkaHandler())
	m.Martini.Use(postgresSession.NewPostgresHandler())
	m.Martini.Use(redisSession.NewRedisHandler())

	log.SetOutput(ioutil.Discard)

	m.Post("/token", func(r render.Render) {
		r.JSON(200, map[string]string{
			"access_token": "ABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDE",
			"token_type":   "bearer",
		})

	})

	// /api/Device/DeviceHeartBeatDetailsNormalized/?id=58e665cd31eb78000738fedd&lastNMilliSeconds=30000&normalizeOnMilliSeconds=250&UniqueDeviceId=HT37WW902113
	m.Get("/api/Device/DeviceHeartBeatDetailsNormalized/",
		func(r render.Render,
			req *http.Request,
			db *mgo.Database) {

			lastRecordID := req.URL.Query().Get("id")
			lastNMilliSeconds, _ := strconv.Atoi(req.URL.Query().Get("lastNMilliSeconds"))
			normalizeOnMilliSeconds, _ := strconv.Atoi(req.URL.Query().Get("normalizeOnMilliSeconds"))
			uniqueDeviceID := req.URL.Query().Get("UniqueDeviceId")

			heartBeatDetails := fetchAllHeartBeatDetailsNormalized(db, lastRecordID, lastNMilliSeconds, normalizeOnMilliSeconds, uniqueDeviceID)

			apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}

			apiResponse.Result = heartBeatDetails

			r.JSON(apiResponse.StatusCode, apiResponse)
		})

	// /api/device/DeviceHeartBeatDetails/?UniqueDeviceId=${UNIQUE_DEVICE_ID}&lastNMilliSeconds=6000
	m.Get("/api/Device/DeviceHeartBeatDetails/",
		func(r render.Render,
			req *http.Request,
			db *mgo.Database) {

			lastRecordID := req.URL.Query().Get("id")
			deviceID := req.URL.Query().Get("UniqueDeviceId")

			lastNMilliSeconds, err := strconv.Atoi(req.URL.Query().Get("lastNMilliSeconds"))

			if err != nil {
				r.JSON(400, "Expecting lastNMilliSeconds parameter in the query string")
				return
			}

			heartBeatDetails := fetchAllHeartBeatDetails(db, lastRecordID, deviceID, lastNMilliSeconds)

			r.JSON(200, heartBeatDetails)
		})

	m.Get("/api/Device/DeviceRpmAnalytics",
		func(r render.Render,
			req *http.Request,
			db *mgo.Database) {

			apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}

			deviceRpmAnalytics, err := fetchDeviceRpmAnalytics(db)

			if err != nil {
				apiResponse.Code = "0001"
				apiResponse.StatusCode = 400
				apiResponse.Message = err.Error()

				r.JSON(apiResponse.StatusCode, apiResponse)
				return
			}

			apiResponse.Result = deviceRpmAnalytics

			r.JSON(200, apiResponse)
		})

	// Define the "GET /api/device/list" route.
	m.Get("/api/Device/DeviceList", func(r render.Render, postgres *PostgresSession) {
		apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}

		devices, err := postgres.FetchAlDevices()

		if err != nil {
			// insert failed, 400 Bad Request
			apiResponse.Code = "0001"
			apiResponse.StatusCode = 400
			apiResponse.Message = err.Error()

			r.JSON(apiResponse.StatusCode, apiResponse)
			return
		}

		apiResponse.Result = devices

		r.JSON(apiResponse.StatusCode, apiResponse)
	})

	// Define the "POST /api/device/test" route.
	m.Post("/api/device/test", binding.Json(HeartBeat{}),
		func(heartBeat HeartBeat,
			r render.Render,
			db *mgo.Database,
			kafka *KafkaSession) {

			apiResponse := utility.APIResponse{}

			apiResponse.Code = "0000"

			r.JSON(200, apiResponse)
		})

	// Define the "POST /api/device/Register" route.
	m.Post("/api/device/RegisterDevice", binding.Json(Device{}),
		func(device Device,
			r render.Render,
			db *mgo.Database,
			postgres *PostgresSession,
			kafka *KafkaSession) {

			err := postgres.RegisterDevice(device)

			// TODO: gseng - fill this function in
			apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}

			if err != nil {
				apiResponse.Result = err.Error()
				apiResponse.Message = "ERROR"
			}

			r.JSON(200, apiResponse)
		})

	// Define the "POST /api/device/SetHeartBeat" route.
	m.Post("/api/device/SetHeartBeat", binding.Json(HeartBeat{}),
		func(heartBeat HeartBeat,
			r render.Render,
			db *mgo.Database,
			postgres *PostgresSession,
			redis *RedisSession,
			kafka *KafkaSession) {

			apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}

			if heartBeat.valid() {
				// device data is valid, now check whether we know this device or not

				// first check the cache
				if !redis.KeyExists(heartBeat.UniqueDeviceID) {
					// if it is not in the cache, query from the db now
					if !postgres.DoesDeviceExist(heartBeat.UniqueDeviceID) {
						// this device is not in the db, send an error
						apiResponse.Code = "0001"
						apiResponse.StatusCode = 400
						apiResponse.Message = fmt.Sprintf("No device exists with id: %s", heartBeat.UniqueDeviceID)
						r.JSON(apiResponse.StatusCode, apiResponse)
						return
					} else {
						// add device to the cache
						redis.SetKey(heartBeat.UniqueDeviceID, "exist")
					}
				}

				// insert into Kafka

				heartBeat.HeartBeatOn = time.Now()
				buff, err := json.Marshal(heartBeat)

				if err != nil {
					apiResponse.Code = "0001"
					apiResponse.StatusCode = 400
					apiResponse.Message = err.Error()
					r.JSON(apiResponse.StatusCode, apiResponse)
					return
				}

				// TODO: gseng - get topic name from Config
				err = kafka.ProduceMessage(string(buff), "SampleApp")
				if err == nil {
					// insert successful, return code should be 201 Created
					// but for compatibility we retur 200
					// r.JSON(200, apiResponse)
					r.JSON(200, apiResponse)
					return
				} else {
					// insert failed, 400 Bad Request
					apiResponse.Code = "0001"
					apiResponse.StatusCode = 400
					apiResponse.Message = err.Error()

					r.JSON(apiResponse.StatusCode, apiResponse)
					return
				}
			} else {
				// heartBeat is invalid, 400 Bad Request
				apiResponse.Code = "0001"
				apiResponse.StatusCode = 400
				apiResponse.Message = "Not a valid device heartbeat"

				r.JSON(apiResponse.StatusCode, apiResponse)
				return
			}
		})

	// Return the server. Call Run() on the server to
	// begin listening for HTTP requests.
	return m
}
