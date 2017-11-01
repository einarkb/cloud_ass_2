package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"../types"
)

const (
	hosts = "ds229435.mlab.com:29435"
	databaseName = "currencydb"
	userName = "test"
	password = "test"
)

func createSession() (*mgo.Session, error) {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  20 * time.Second,
		Database: databaseName,
		Username: userName,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Fatal("Failed to create database connection")
		return nil, err
	}

	return session, nil
}

func InsertCurrencyTick(currencyTick types.CurrencyData) {
	session, err := createSession()
	if err != nil {
		return
	}
	defer session.Close()

	err2 := session.DB("currencydb").C("tick").Insert(currencyTick)
	if err2 != nil {
		log.Fatal("Error inserting currencytick to db")
	}
}

// returns the "r" last ticks from index "i" (1 = last)
// 1, 1 will return the newest tick, 1, 3 the 3 newest ticks
func GetCurrencies(i int, r int) []types.CurrencyData {
	if r <= 0 || i <= 0 {
		log.Fatal("Error inserting currencytick to db")
		return nil
	}

	session, err := createSession()
	if err != nil {
		return nil
	}
	defer session.Close()

	collection := session.DB(databaseName).C("tick")
	dbSize, err2 := collection.Count()
	if err2 != nil {
		return nil
	}

	var data[] types.CurrencyData
	dif := r - dbSize
	if dif < 0 {
		dif = 0
	}
	collection.Find(nil).Skip(dbSize - i - r + 1).Limit(r).All(&data)
	return data
}
