package main

import (
	"net/http"
	"fmt"
	"strings"
	"../db"
	//"../fixer"
	"../types"
)

func handlerLatest(w http.ResponseWriter, r *http.Request) {
	langs := strings.Split(strings.Split(r.URL.Path, "latest/")[1], "/")

	data := db.GetCurrencies(1, 1)
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

func handlerAverage(w http.ResponseWriter, r *http.Request) {
	langs := strings.Split(strings.Split(r.URL.Path, "average/")[1], "/")

	data := db.GetCurrencies(1, 7)
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
	http.ListenAndServe("127.0.0.1:8081", nil)
}