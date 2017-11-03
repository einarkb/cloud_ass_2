package fixer

import (
	"testing"
	"github.com/einarkb/cloud_ass_2/db"
	"net/http"
	"github.com/einarkb/cloud_ass_2/types"
	"encoding/json"
)

func TestDBCurrency(t *testing.T) {
	db.ClearTestCollection()
	resp, err := http.Get("http://api.fixer.io/latest?base=EUR")
	if err != nil {
		t.Error("could not get data from fixer")
	}
	payload := types.CurrencyData{}
	json.NewDecoder(resp.Body).Decode(&payload)

	db.InsertCurrencyTick(payload, "test")
	res := db.GetCurrencies(1, 1, "test")

	if res[0].Rates["NOK"] != payload.Rates["NOK"] {
		t.Error("currency did not match before and after db insertion and retreival")
	}
	db.ClearTestCollection()
}

func Test_fetchDataFromFixer(t *testing.T) {
	db.ClearTestCollection()
}