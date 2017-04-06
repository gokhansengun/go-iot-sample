package iot

// RpmAnalyticsHeartBeat is composed of RPM values
type RpmAnalyticsHeartBeat struct {
	UniqueDeviceID            string `json:"uniqueDeviceId" bson:"uniqueDeviceId"`
	NumberOfRequestPerMinutes int    `json:"numberOfRequestPerMinutes" bson:"numberOfRequestPerMinutes"`
}

// RpmAnalyticsHeartBeats is array of RpmAnalyticsHeartBeat objects
type RpmAnalyticsHeartBeats []RpmAnalyticsHeartBeat
