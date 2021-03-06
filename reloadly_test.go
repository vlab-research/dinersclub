package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/nandanrao/go-reloadly/reloadly"
	"github.com/stretchr/testify/assert"
)

func TestReloadlyResultsOnErrorIfBadDetails(t *testing.T) {
	provider := &ReloadlyProvider{&ReloadlyConfig{}, &reloadly.Service{
		Client: &http.Client{},
	}}

	jm := json.RawMessage([]byte(`{"foo": "bar"}`))

	ts := JSTimestamp(time.Now().UTC())

	pe := &PaymentEvent{
		Userid:    "foo",
		Pageid:    "page",
		Timestamp: &ts,
		Provider:  "reloadly",
		Details:   &jm,
	}

	res, err := provider.Payout(pe)
	assert.Nil(t, err)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "INVALID_PAYMENT_DETAILS", res.Error.Code)
	assert.Equal(t, "payment:reloadly", res.Type)
	assert.Equal(t, false, res.Success)
}

func TestReloadlyReportsAPIErrorsInResult(t *testing.T) {
	provider := &ReloadlyProvider{&ReloadlyConfig{}, &reloadly.Service{
		Client: TestClient(404, `{"errorCode": "FOOBAR", "message": "Sorry"}`, nil),
	}}

	jm := json.RawMessage([]byte(`{"number": "+123", "amount": 2.5, "country": "IN", "id": "id"}`))
	ts := JSTimestamp(time.Now().UTC())
	pe := &PaymentEvent{
		Userid:    "foo",
		Pageid:    "page",
		Timestamp: &ts,
		Provider:  "reloadly",
		Details:   &jm,
	}

	res, err := provider.Payout(pe)
	assert.Nil(t, err)
	assert.NotNil(t, res.Error)
	assert.Equal(t, "FOOBAR", res.Error.Code)
	assert.Equal(t, "Sorry", res.Error.Message)
	assert.Equal(t, &jm, res.Error.PaymentDetails)
	assert.Equal(t, "id", res.ID)
	assert.Equal(t, "payment:reloadly", res.Type)
	assert.Equal(t, false, res.Success)
}

func TestReloadlyReportsSuccessResult(t *testing.T) {
	provider := &ReloadlyProvider{&ReloadlyConfig{}, &reloadly.Service{
		Client: TestClient(200, `{"suggestedAmountsMap":{"2.5": 2.5},"transactionDate":"2020-09-19 12:53:22","transactionId": 567}`, nil),
	}}

	jm := json.RawMessage([]byte(`{"number": "+123", "amount": 2.5, "country": "IN"}`))
	ts := JSTimestamp(time.Now().UTC())
	pe := &PaymentEvent{
		Userid:    "foo",
		Pageid:    "page",
		Timestamp: &ts,
		Provider:  "reloadly",
		Details:   &jm,
	}

	res, err := provider.Payout(pe)
	assert.Nil(t, err)
	assert.Nil(t, res.Error)
	assert.Equal(t, "payment:reloadly", res.Type)
	assert.Equal(t, true, res.Success)
	assert.Equal(t, &jm, res.PaymentDetails)
}
