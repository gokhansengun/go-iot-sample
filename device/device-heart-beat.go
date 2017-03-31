package device

import (
	"github.com/gokhansengun/go-iot-sample/utility"
)

// HeartBeat is composed of a first name, last name,
// email, age, and short message. When represented in
// JSON, ditch TitleCase for snake_case.
type HeartBeat struct {
	UniqueDeviceID string `json:"UniqueDeviceId"`

	AccelerometerX float32
	AccelerometerY float32
	AccelerometerZ float32

	GyroscopeX float32
	GyroscopeY float32
	GyroscopeZ float32

	MagnetometerX float32
	MagnetometerY float32
	MagnetometerZ float32

	Compass   float32
	Latitude  float32
	Longitude float32

	HeartBeatOn utility.CustomTime
}

// All fields must exist and valid
func (heartBeat *HeartBeat) valid() bool {
	// TODO: gseng - no checks for the time being
	return true
}
