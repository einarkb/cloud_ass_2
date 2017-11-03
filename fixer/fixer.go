package fixer

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	"time"
	"github.com/einarkb/clodsud_ass_2/db"
	"github.com/einarkb/cloud_ass_2/types"
)

//fetchDataFromFixer reuqtest data from the fixer and stores it in the database
func fetchDataFromFixer() {
	resp, err := http.Get("http://api.fixer.io/latest?base=EUR")
	if err != nil {
		log.Fatal(err)
		return
	}

	payload := types.CurrencyData{}
	json.NewDecoder(resp.Body).Decode(&payload)

	db.InsertCurrencyTick(payload, "tick")

	fmt.Println("tick")
}

// startTicker creates a loop that will call fetchDataFromFixer once every 24 hours
func startTicker() {
	fetchDataFromFixer()
	ticker := time.NewTicker(time.Second)
	for {
		fetchDataFromFixer()
		<-ticker.C
	}
}

// Start creates a thread that runs the fixer
func Start() {
	go startTicker()
}