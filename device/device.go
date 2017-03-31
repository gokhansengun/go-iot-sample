package device

import (
	"labix.org/v2/mgo"
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
