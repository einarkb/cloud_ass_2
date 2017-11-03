package types

import "gopkg.in/mgo.v2/bson"

// CurrencyData stores one tick of info about currencies
type CurrencyData struct {
	Base  string
	Date  string
	Rates map[string]float64
}

// WebhookPayload keeps info of a webhook
type WebhookPayload struct {
	ID              bson.ObjectId `json:"id" bson:"_id"`
	WebhookURL      string        `json:"webhookURL" bson:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency" bson:"targetCurrency"`
	MinTriggerValue float64       `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64       `json:"maxTriggerValue" bson:"maxTriggerValue"`
	CurrentRate     float64       `json:"currentRate" bson:"currentRate"`
}

// WebhookPayload keeps info for invoking a webhook
type WebhookInvokePayload struct {
	BaseCurrency    string  `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency" bson:"targetCurrency"`
	CurrentRate     float64 `json:"currentRate" bson:"currentRate"`
	MinTriggerValue float64 `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float64 `json:"maxTriggerValue" bson:"maxTriggerValue"`
}