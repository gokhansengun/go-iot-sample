package device

import (
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Device is composed of some variables
type Device struct {
	UniqueDeviceID string `json:"UniqueDeviceId"`
	DeviceType     string
}

// All fields must exist and valid
func (device *Device) valid() bool {
	// TODO: gseng - no checks for the time being
	return true
}

// retrieve all device details without any condition
func fetchAllDevices(db *mgo.Database) []Device {
	devices := []Device{}
	err := db.C("Device").Find(nil).All(&devices)
	if err != nil {
		panic(err)
	}

	return devices
}

func fetchAllHeartBeatDetails(db *mgo.Database, deviceID string, lastNMilliSeconds int) []HeartBeat {
	heartBeats := []HeartBeat{}

	now := time.Now()
	duration := time.Duration(-lastNMilliSeconds) * time.Millisecond
	cutoffDate := now.Add(duration)

	query := bson.M{"$and": []bson.M{bson.M{"UniqueDeviceId": deviceID}, bson.M{"HeartBeatOn": bson.M{"$gt": cutoffDate}}}}

	if deviceID == "" {
		query = bson.M{"HeartBeatOn": bson.M{"$gt": cutoffDate}}
	}

	err := db.C("DeviceDetail").Find(query).All(&heartBeats)
	if err != nil {
		panic(err)
	}

	return heartBeats
}
