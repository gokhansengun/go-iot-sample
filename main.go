package main

import (
	"github.com/gokhansengun/go-iot-sample/device"
)

/*
Create a new MongoDB session, using a database
named "SampleApp". Create a new server using
that session, then begin listening for HTTP requests.
*/
func main() {

	appConf := device.NewConfig("conf.yaml")

	// TODO: gseng - check whether below is the database name, if so change it
	mongoDbSession := device.NewMongoDbSession("SampleApp", appConf.GetMongoDbConnStr())
	kafkaSession := device.NewKafkaSession(appConf.GetKafkaConnStr())
	server := device.NewServer(mongoDbSession, kafkaSession)
	server.Martini.Run()
}
