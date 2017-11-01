package types

type CurrencyData struct {
	Base  string
	Date  string
	Rates map[string]float64
}
