package main

import (
	"github.com/globalsign/mgo"
)

// CreateSession returns a session to our database
// it takes a host string as an argument.
func CreateSession(host string) (*mgo.Session, error) {
	session, err := mgo.Dial(host)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}
