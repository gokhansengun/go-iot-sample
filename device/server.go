package device

import (
	"encoding/json"
	"io/ioutil"
	"log"

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
func NewServer(session *DatabaseSession, kafkaSession *KafkaSession) *MartiniServer {
	m := initMartiniServer()
	m.Martini.Use(render.Renderer(render.Options{
		IndentJSON: true,
	}))
	m.Martini.Use(session.NewMongoDbDatabase())
	m.Martini.Use(kafkaSession.NewKafkaSyncQueue())

	log.SetOutput(ioutil.Discard)

	m.Post("/token", func(r render.Render) {
		r.JSON(200, map[string]string{
			"access_token": "ABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDEABCDE",
			"token_type":   "bearer",
		})

	})

	// Define the "GET /api/device/list" route.
	m.Get("/api/device/list", func(r render.Render, db *mgo.Database, kafka *KafkaSession) {
		r.JSON(200, fetchAllDevices(db))
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
			kafka *KafkaSession) {

			// TODO: gseng - fill this function in
			apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}
			r.JSON(200, apiResponse)
		})

	// Define the "POST /api/device/SetHeartBeat" route.
	m.Post("/api/device/SetHeartBeat", binding.Json(HeartBeat{}),
		func(heartBeat HeartBeat,
			r render.Render,
			db *mgo.Database,
			kafka *KafkaSession) {

			apiResponse := utility.APIResponse{StatusCode: 200, Code: "0000", Message: "OK"}

			if heartBeat.valid() {
				// device is valid, insert into Kafka
				buff, err := json.Marshal(heartBeat)

				if err != nil {
					apiResponse.Code = "0001"
					apiResponse.StatusCode = 400
					apiResponse.Message = err.Error()
					r.JSON(apiResponse.StatusCode, apiResponse)
				}

				// TODO: gseng - get topic name from Config
				err = kafka.ProduceMessage(string(buff), "SampleApp")
				if err == nil {
					// insert successful, return code should be 201 Created
					// but for compatibility we retur 200
					// r.JSON(200, apiResponse)
					r.JSON(200, apiResponse)
				} else {
					// insert failed, 400 Bad Request
					apiResponse.Code = "0001"
					apiResponse.StatusCode = 400
					apiResponse.Message = err.Error()

					r.JSON(apiResponse.StatusCode, apiResponse)
				}
			} else {
				// heartBeat is invalid, 400 Bad Request
				apiResponse.Code = "0001"
				apiResponse.StatusCode = 400
				apiResponse.Message = "Not a valid device heartbeat"

				r.JSON(apiResponse.StatusCode, apiResponse)
			}
		})

	// Return the server. Call Run() on the server to
	// begin listening for HTTP requests.
	return m
}