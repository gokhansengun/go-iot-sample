package main

import (
	"github.com/gokhansengun/go-iot-sample/iot"
)

/*
Create a new MongoDB session, using a database
named "SampleApp". Create a new server using
that session, then begin listening for HTTP requests.
*/
func main() {

	appConf := iot.NewConfig("conf.yaml")

	// TODO: gseng - check whether below is the database name, if so change it
	mongoDbSession := iot.NewMongoDbSession("SampleApp", appConf.GetMongoDbConnStr())
	kafkaSession := iot.NewKafkaSession(appConf.GetKafkaConnStr())
	postgresSession := iot.NewPostgresSession(appConf.GetPostgresConnStr())

	server := iot.NewServer(mongoDbSession, kafkaSession, postgresSession)
	server.Martini.Run()
}
