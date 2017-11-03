package main

import (
	"testing"
	//"../types"
	"net/http"
)

func Test_handleWebhookPost(t *testing.T) {
	//webhook := types.WebhookPayload{}
	id := "59fa692550ad253394b591fa"

	_, err := http.NewRequest("GET", "/" + id, nil)
	if err != nil {
		t.Error("respone from GET request is nil")
	}
}