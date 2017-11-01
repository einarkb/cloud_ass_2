package db

import (
	"gopkg.in/mgo.v2"
	"log"
)

func DbConnect() *mgo.Session {
	mongoDbConnect := "mongodb://test:test@ds229435.mlab.com:29435/currencydb"
	session, err := mgo.Dial(mongoDbConnect)
	if err != nil {
		log.Fatal("mongodb failed @ ", mongoDbConnect, err.Error())
		return nil
	}
	return session
}
