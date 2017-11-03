package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
	"../types"
	"gopkg.in/mgo.v2/bson"
)

const (
	hosts = "ds229435.mlab.com:29435"
	databaseName = "currencydb"
	userName = "test"
	password = "test"
)

//createSession connects to the database and returns the session
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
	session.SetMode(mgo.Monotonic, true)

	return session, nil
}

//InsertCurrencyTick inserts a tick from the fixer into the database
func InsertCurrencyTick(currencyTick types.CurrencyData, col string) {
	session, err := createSession()
	if err != nil {
		return
	}
	defer session.Close()

	err2 := session.DB("currencydb").C(col).Insert(currencyTick)
	if err2 != nil {
		log.Fatal("Error inserting currencytick to db")
	}
}

// GetCurrencies returns the "r" last ticks from index "i" (1 = last)
// 1, 1 will return the newest tick, 1, 3 the 3 newest ticks
func GetCurrencies(i int, r int, col string) []types.CurrencyData {
	if r <= 0 || i <= 0 {
		log.Fatal("Error inserting currencytick to db")
		return nil
	}

	session, err := createSession()
	if err != nil {
		return nil
	}
	defer session.Close()

	collection := session.DB(databaseName).C(col)
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

//InsertWebhook inserts a webhook into the database
func InsertWebhook(webhook types.WebhookPayload) {
	session, err := createSession()
	if err != nil {
		return
	}
	defer session.Close()

	session.DB(databaseName).C("webhooks").Insert(&webhook)
}

//GetWebhook gets webook with id "s" from databse and returns it
func GetWebhook(s string) (types.WebhookPayload, error) {
	payload := types.WebhookPayload{}

	session, err := createSession()
	if err != nil {
		return payload, err
	}
	defer session.Close()

	err1 := session.DB(databaseName).C("webhooks").FindId(bson.ObjectIdHex(s)).One(&payload)
	if err1 != nil {
		return payload, err
	}

	return payload, nil
}

//DeleteWebhook deletes webook with id "s" form databse
func DeleteWebhook(s string) error {
	session, err := createSession()
	if err != nil {
		return err
	}
	defer session.Close()

	err1 := session.DB(databaseName).C("webhooks").RemoveId(bson.ObjectIdHex(s))
	if err1 != nil {
		return err1
	}

	return nil
}

//ClearTestCollection clears the test database
func ClearTestCollection() {
	session, err := createSession()
	if err != nil {
		return
	}
	defer session.Close()
	session.DB(databaseName).C("test").RemoveAll(nil)
}