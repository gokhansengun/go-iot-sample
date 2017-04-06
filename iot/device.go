package iot

import (
	"sort"
	"time"

	// Using blank import for Postgres Driver
	_ "github.com/lib/pq"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// Device is composed of some variables
type Device struct {
	UniqueDeviceID string `json:"uniqueDeviceId"`
	DeviceType     int    `json:"deviceType"`
}

// All fields must exist and valid
func (device *Device) valid() bool {
	// TODO: gseng - no checks for the time being
	return true
}

func fetchAllHeartBeatDetails(db *mgo.Database, lastRecordID string, deviceID string, lastNMilliSeconds int) []HeartBeat {
	heartBeats := []HeartBeat{}

	now := time.Now()
	duration := time.Duration(-lastNMilliSeconds) * time.Millisecond
	cutoffDate := now.Add(duration)

	query := bson.M{"$and": []bson.M{bson.M{"UniqueDeviceId": deviceID}, bson.M{"HeartBeatOn": bson.M{"$gt": cutoffDate}}}}

	if lastRecordID != "" {
		query = bson.M{"$and": []bson.M{bson.M{"UniqueDeviceId": deviceID}, bson.M{"_id": bson.M{"$gt": bson.ObjectIdHex(lastRecordID)}}, bson.M{"HeartBeatOn": bson.M{"$gt": cutoffDate}}}}
	}

	err := db.C("DeviceDetail").Find(query).All(&heartBeats)
	if err != nil {
		panic(err)
	}

	return heartBeats
}

func fetchAllHeartBeatDetailsNormalized(db *mgo.Database, lastRecordID string, lastNMilliSeconds int, normalizeOnMilliSeconds int, uniqueDeviceID string) []HeartBeat {
	// first retrieve the heartbeat details for the lastNMilliSeconds
	heartBeatsLastNMilliSeconds := fetchAllHeartBeatDetails(db, lastRecordID, uniqueDeviceID, lastNMilliSeconds)

	// now normalize and group on each normalizeOnMilliSeconds (example 100 milliseconds)
	normalizationMap := make(map[int64][]HeartBeat)

	normalizeOnNanoSeconds := int64(normalizeOnMilliSeconds * 1000000)

	for _, heartBeat := range heartBeatsLastNMilliSeconds {
		inUnixNanoseconds := heartBeat.HeartBeatOn.UnixNano()
		inUnixNanoseconds = inUnixNanoseconds - (inUnixNanoseconds % normalizeOnNanoSeconds)

		newHeartBeat := heartBeat
		newHeartBeat.HeartBeatOn = time.Unix(0, inUnixNanoseconds)

		// check whether this entry exists
		heartBeatGroup, ok := normalizationMap[inUnixNanoseconds]

		if !ok {
			// add the group to the map
			heartBeatGroup = []HeartBeat{}
		}

		heartBeatGroup = append(heartBeatGroup, newHeartBeat)
		normalizationMap[inUnixNanoseconds] = heartBeatGroup
	}

	normalizedHeartBeats := HeartBeats{}

	// Normalization and grouping complete, create only one entry for each group
	for _, group := range normalizationMap {
		maxID := group[0].ID

		for _, heartBeat := range group {
			if heartBeat.ID > maxID {
				maxID = heartBeat.ID
			}
		}

		normalizedHeartBeat := group[0]
		normalizedHeartBeat.ID = maxID

		normalizedHeartBeats = append(normalizedHeartBeats, normalizedHeartBeat)
	}

	sort.Sort(normalizedHeartBeats)

	return normalizedHeartBeats
}

func fetchDeviceRpmAnalytics(db *mgo.Database) (RpmAnalyticsHeartBeats, error) {
	now := time.Now()
	duration := time.Duration(-1) * time.Minute
	cutoffDate := now.Add(duration)

	query := bson.M{"HeartBeatOn": bson.M{"$gt": cutoffDate}}

	count, err := db.C("DeviceDetail").Find(query).Count()

	if err != nil {
		return nil, err
	}

	rpmHeartBeat := RpmAnalyticsHeartBeat{UniqueDeviceID: "Dummy", NumberOfRequestPerMinutes: count}

	rpmHeartBeats := RpmAnalyticsHeartBeats{}
	rpmHeartBeats = append(rpmHeartBeats, rpmHeartBeat)

	return rpmHeartBeats, nil
}
