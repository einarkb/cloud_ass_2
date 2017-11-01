package fixer

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"../db"
)

type FixerPayload struct {
	Base  string
	Date  string
	Rates map[string]float64
}

func fetchDataFromFixer() {
	resp, err := http.Get("http://api.fixer.io/latest?base=EUR")
	if err != nil {
		log.Fatal(err)
		return
	}

	payload := FixerPayload{}
	json.NewDecoder(resp.Body).Decode(&payload)

	session := db.DbConnect()
	if session == nil {
		return
	}
	defer session.Close()

	err2 := session.DB("currencydb").C("tick").Insert(payload)
	if err2 != nil {
		log.Fatal("Error on session.DB(", "currencydb", ").C(", "tick", ").Insert(<Payload>)", err2.Error())
	}

	fmt.Println("tick")
}

func startTicker() {
	ticker := time.NewTicker(time.Second)
	for {
		fetchDataFromFixer()
		<-ticker.C
	}
}

func Start() {
	go startTicker()
}