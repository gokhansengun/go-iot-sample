package iot

// RpmAnalyticsHeartBeat is composed of RPM values
type RpmAnalyticsHeartBeat struct {
	UniqueDeviceID            string `bson:"uniqueDeviceId"`
	NumberOfRequestPerMinutes int    `bson:"numberOfRequestPerMinutes"`
}

// RpmAnalyticsHeartBeats is array of RpmAnalyticsHeartBeat objects
type RpmAnalyticsHeartBeats []RpmAnalyticsHeartBeat
