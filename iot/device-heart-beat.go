package iot

import (
	"time"

	"labix.org/v2/mgo/bson"
)

// HeartBeat is composed of a first name, last name,
// email, age, and short message. When represented in
// JSON, ditch TitleCase for snake_case.
type HeartBeat struct {
	ID bson.ObjectId `json:"-" bson:"_id,omitempty"`

	UniqueDeviceID string `bson:"UniqueDeviceId"`

	AccelerometerX float64 `bson:"AccelerometerX"`
	AccelerometerY float64 `bson:"AccelerometerY"`
	AccelerometerZ float64 `bson:"AccelerometerZ"`

	GyroscopeX float64 `bson:"GyroscopeX"`
	GyroscopeY float64 `bson:"GyroscopeY"`
	GyroscopeZ float64 `bson:"GyroscopeZ"`

	MagnetometerX float64 `bson:"MagnetometerX"`
	MagnetometerY float64 `bson:"MagnetometerY"`
	MagnetometerZ float64 `bson:"MagnetometerZ"`

	Compass   float64 `bson:"Compass"`
	Latitude  float64 `bson:"Latitude"`
	Longitude float64 `bson:"Longitude"`

	HeartBeatOn time.Time `bson:"HeartBeatOn"`
}

// HeartBeats is the array of HeartBeat objects
type HeartBeats []HeartBeat

func (slice HeartBeats) Len() int {
	return len(slice)
}

func (slice HeartBeats) Less(i, j int) bool {
	return slice[i].HeartBeatOn.Before(slice[j].HeartBeatOn)
}

func (slice HeartBeats) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// All fields must exist and valid
func (heartBeat *HeartBeat) valid() bool {
	// TODO: gseng - no checks for the time being
	return true
}
