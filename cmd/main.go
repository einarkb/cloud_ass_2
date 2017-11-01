package main


import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
)

type FixerPayload struct {
	Base string
	Date string
	Rates map[string]float64
}

func dbConnect() *mgo.Session {
	mongoDbConnect := "mongodb://test:test@ds229435.mlab.com:29435/currencydb"
	session, err := mgo.Dial(mongoDbConnect)
	if err != nil {
		log.Fatal("mongodb failed @ ", mongoDbConnect, err.Error())
		return nil
	}
	return session
}

func fetchDataFromFixer() {
	resp, err := http.Get("http://api.fixer.io/latest?base=EUR")
	if err != nil {
		log.Fatal(err)
		return
	}

	payload := FixerPayload{}
	json.NewDecoder(resp.Body).Decode(&payload)

	session := dbConnect()
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

func handlerFunc(w http.ResponseWriter, r *http.Request) {

}

func handlerLatest(w http.ResponseWriter, r *http.Request) {
	langs := strings.Split(strings.Split(r.URL.Path, "latest/")[1], "/")

	session := dbConnect()
	if session == nil {
		return
	}
	defer session.Close()

	collection := session.DB("currencydb").C("tick")
	dbSize, err := collection.Count()
	if err != nil {
		return
	}
	var data FixerPayload
	collection.Find(nil).Skip(dbSize - 1).One(&data)

	valid := isLanguageInputValid(data, langs)
	if !valid {
		fmt.Fprint(w, "Invalid currencies")
		return
	}

	if langs[0] == "EUR" {
		fmt.Fprint(w, data.Rates[langs[1]])
	} else if langs[1] == "EUR" {
		fmt.Fprint(w,  1 / data.Rates[langs[0]])
	}else {
		fmt.Fprint(w, data.Rates[langs[1]] / data.Rates[langs[0]])
	}

}

func handlerAverage(w http.ResponseWriter, r *http.Request) {
	langs := strings.Split(strings.Split(r.URL.Path, "average/")[1], "/")

	session := dbConnect()
	if session == nil {
		return
	}
	defer session.Close()

	collection := session.DB("currencydb").C("tick")
	dbSize, err := collection.Count()
	if err != nil {
		return
	}
	var data[] FixerPayload
	dif := 7 - dbSize
	if dif < 0 {
		dif = 0
	}
	collection.Find(nil).Skip(dbSize - 7 + dif).All(&data)
	valid := isLanguageInputValid(data[0], langs)
	if !valid {
		fmt.Fprint(w, "Invalid currencies")
		return
	}

	fmt.Fprint(w, len(data))
	var lang0Avg, lang1Avg float64
	for i := 0; i < len(data); i++ {
		lang0Avg += data[i].Rates[langs[0]]
		lang1Avg += data[i].Rates[langs[1]]
	}
	lang0Avg /= 7 - float64(dif)
	lang1Avg /= 7 - float64(dif)

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
func isLanguageInputValid(data FixerPayload, langs[] string) bool {
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
	//go startTicker()

	http.HandleFunc("/", handlerFunc)
	http.HandleFunc("/latest/", handlerLatest)
	http.HandleFunc("/average/", handlerAverage)
	http.ListenAndServe("127.0.0.1:8081", nil)
}