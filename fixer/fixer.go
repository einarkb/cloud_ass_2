package fixer

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"../db"
	"../types"
)

func fetchDataFromFixer() {
	resp, err := http.Get("http://api.fixer.io/latest?base=EUR")
	if err != nil {
		log.Fatal(err)
		return
	}

	payload := types.CurrencyData{}
	json.NewDecoder(resp.Body).Decode(&payload)

	db.InsertCurrencyTick(payload)

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