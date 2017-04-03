package iot

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
)

// DatabaseSession embed *mgo.Session and store the database name.
type DatabaseSession struct {
	*mgo.Session
	databaseName string
}

// NewMongoDbSession connects to the local MongoDB and set up the database.
func NewMongoDbSession(name string, connStr string) *DatabaseSession {
	session, err := mgo.Dial(connStr)
	if err != nil {
		panic(err)
	}

	return &DatabaseSession{session, name}
}

// NewMongoDbDatabase - Martini lets you inject parameters for routing handlers
// by using `context.Map()`. I'll pass each route handler
// a instance of a *mgo.Database, so they can retrieve
// and insert device heartbeats to and from that database.
// For more information, check out:
// http://blog.gopheracademy.com/day-11-martini
func (session *DatabaseSession) NewMongoDbDatabase() martini.Handler {
	return func(context martini.Context) {
		s := session.Clone()
		context.Map(s.DB(session.databaseName))
		defer s.Close()
		context.Next()
	}
}
