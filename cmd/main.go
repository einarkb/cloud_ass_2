package main

import (
	"net/http"
	"fmt"
	"strings"
	"github.com/einarkb/cloud_ass_2/db"
	"github.com/einarkb/cloud_ass_2/types"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"os"
)

// Prints the newest rate between the 2 specified currencies.
// Supports any currency to any other. feks DKK to SEK
func handlerLatest(w http.ResponseWriter, r *http.Request) {
	langs := strings.Split(strings.Split(r.URL.Path, "latest/")[1], "/")

	data := db.GetCurrencies(1, 1, "tick")
	if data == nil {
		return
	}

	valid := isLanguageInputValid(data[0], langs)
	if !valid {
		fmt.Fprint(w, "Invalid currencies")
		return
	}

	if langs[0] == "EUR" {
		fmt.Fprint(w, data[0].Rates[langs[1]])
	} else if langs[1] == "EUR" {
		fmt.Fprint(w,  1 / data[0].Rates[langs[0]])
	}else {
		fmt.Fprint(w, data[0].Rates[langs[1]] / data[0].Rates[langs[0]])
	}
}

// handlerAverage prints the average rate between the 2 specified currencies for the last 7 days.
// Supports any currency to any other. feks DKK to SEK
func handlerAverage(w http.ResponseWriter, r *http.Request) {
	langs := strings.Split(strings.Split(r.URL.Path, "average/")[1], "/")

	data := db.GetCurrencies(1, 7, "tick")
	if data == nil {
		return
	}

	valid := isLanguageInputValid(data[0], langs)
	if !valid {
		fmt.Fprint(w, "Invalid currencies")
		return
	}

	var lang0Avg, lang1Avg float64
	for i := 0; i < len(data); i++ {
		lang0Avg += data[i].Rates[langs[0]]
		lang1Avg += data[i].Rates[langs[1]]
	}
	lang0Avg /= 7 - float64(7 - len(data))
	lang1Avg /= 7 - float64(7 - len(data))

	if langs[0] == "EUR" {
		fmt.Fprint(w, lang1Avg)
	} else if langs[1] == "EUR" {
		fmt.Fprint(w,  1 / lang0Avg)
	}else {
		fmt.Fprint(w, lang1Avg / lang0Avg)
	}
}

// handleWebhookPost adds webhook specified in the post request
func handleWebhookPost(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "text/plain")

	webhook := types.WebhookPayload{}

	decErr := json.NewDecoder(r.Body).Decode(&webhook)

	if decErr != nil {
		http.Error(w, "Error23: ", http.StatusInternalServerError)
		return
	}

	webhook.ID = bson.NewObjectId()
	webhook.CurrentRate = db.GetCurrencies(1, 1, "tick")[0].Rates[webhook.TargetCurrency]
	db.InsertWebhook(webhook)

	fmt.Fprintf(w, webhook.ID.Hex())
}

// handleWebhookGet gets webhook with id specified by s
func handleWebhookGet(s string, w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")

	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Invalid id: ", http.StatusBadRequest)
		return
	}

	webhook, err := db.GetWebhook(s)
	if err != nil {
		http.Error(w, "Could not get the webhook form database: ", http.StatusInternalServerError)
	}

	webhook.CurrentRate = db.GetCurrencies(1, 1, "tick")[0].Rates[webhook.TargetCurrency]

	json.NewEncoder(w).Encode(webhook)
}

// handleWebhookDelete Deleted weebook with id specified by s
func handleWebhookDelete(s string, w http.ResponseWriter, r *http.Request) {
	if bson.IsObjectIdHex(s) == false {
		http.Error(w, "Invalid id: ", http.StatusBadRequest)
		return
	}

	err := db.DeleteWebhook(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// handlerRoot controls what the server does in root depending on the request
func handlerRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleWebhookPost(w, r)

	case "GET":
		handleWebhookGet(strings.Split(r.URL.Path, "/")[1], w, r)

	case "DELETE":
		handleWebhookDelete(strings.Split(r.URL.Path, "/")[1], w, r)

	default:
		http.Error(w, "Request not supported.", http.StatusNotImplemented)
		return
	}
}

// checks if the specified currencies actually exists and are not duplicates.
// if so, ok wil be set to true, otherwise false
func isLanguageInputValid(data types.CurrencyData, langs[] string) bool {
	var ok bool
	for i := 0; i <= 1; i++ {
		_, ok = data.Rates[langs[i]]
		if !ok && langs[i] != data.Base {
			break
		} else {
			ok = true
		}
	}
	if langs[0] == langs[1] {
		ok = false
	}
	if ok {
		return true
	} else {
		return false
	}
}

func main() {
	//fixer.Start()

	http.HandleFunc("/latest/", handlerLatest)
	http.HandleFunc("/average/", handlerAverage)
	http.HandleFunc("/", handlerRoot)
	http.ListenAndServe("8080", nil)
}