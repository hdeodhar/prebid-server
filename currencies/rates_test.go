package currencies_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prebid/prebid-server/currencies"
)

func TestUnMarshallRates(t *testing.T) {

	// Setup:
	testCases := []struct {
		ratesJSON     string
		expectedRates currencies.Rates
		expectsError  bool
	}{
		{
			ratesJSON: `{
				"dataAsOf":"2018-09-12",
				"conversions":{
					"USD":{
						"GBP":0.7662523901
					},
					"GBP":{
						"USD":1.3050530256
					}
				}
			}`,
			expectedRates: currencies.Rates{
				DataAsOf: time.Date(2018, time.September, 12, 0, 0, 0, 0, time.UTC),
				Conversions: map[string]map[string]float64{
					"USD": {
						"GBP": 0.7662523901,
					},
					"GBP": {
						"USD": 1.3050530256,
					},
				},
			},
			expectsError: false,
		},
		{
			ratesJSON: `{
				"dataAsOf":"",
				"conversions":{
					"USD":{
						"GBP":0.7662523901
					},
					"GBP":{
						"USD":1.3050530256
					}
				}
			}`,
			expectedRates: currencies.Rates{
				DataAsOf: time.Time{},
				Conversions: map[string]map[string]float64{
					"USD": {
						"GBP": 0.7662523901,
					},
					"GBP": {
						"USD": 1.3050530256,
					},
				},
			},
			expectsError: false,
		},
		{
			ratesJSON: `{
				"dataAsOf":"blabla",
				"conversions":{
					"USD":{
						"GBP":0.7662523901
					},
					"GBP":{
						"USD":1.3050530256
					}
				}
			}`,
			expectedRates: currencies.Rates{
				DataAsOf: time.Time{},
				Conversions: map[string]map[string]float64{
					"USD": {
						"GBP": 0.7662523901,
					},
					"GBP": {
						"USD": 1.3050530256,
					},
				},
			},
			expectsError: false,
		},
		{
			ratesJSON: `{
				"dataAsOf":"blabla",
				"conversions":{
					"USD":{
						"GBP":0.7662523901,
					},
					"GBP":{
						"USD":1.3050530256,
					}
				}
			}`,
			expectedRates: currencies.Rates{},
			expectsError:  true,
		},
	}

	for _, tc := range testCases {

		// Execute:
		updatedRates := currencies.Rates{}
		err := json.Unmarshal([]byte(tc.ratesJSON), &updatedRates)

		// Verify:
		assert.Equal(t, err != nil, tc.expectsError)
		assert.Equal(t, tc.expectedRates, updatedRates, "Rates weren't the expected ones")
	}
}

func TestGetRate(t *testing.T) {

	// Setup:
	rates := currencies.NewRates(time.Now(), map[string]map[string]float64{
		"USD": {
			"GBP": 0.77208,
		},
		"GBP": {
			"USD": 1.2952,
		},
	})

	testCases := []struct {
		from         string
		to           string
		expectedRate float64
		hasError     bool
	}{
		{from: "USD", to: "GBP", expectedRate: 0.77208, hasError: false},
		{from: "GBP", to: "USD", expectedRate: 1.2952, hasError: false},
		{from: "GBP", to: "EUR", expectedRate: 0, hasError: true},
		{from: "CNY", to: "EUR", expectedRate: 0, hasError: true},
		{from: "", to: "EUR", expectedRate: 0, hasError: true},
		{from: "CNY", to: "", expectedRate: 0, hasError: true},
		{from: "", to: "", expectedRate: 0, hasError: true},
	}

	// Verify:
	for _, tc := range testCases {
		rate, err := rates.GetRate(tc.from, tc.to)

		if tc.hasError {
			assert.NotNil(t, err, "err shouldn't be nil")
			assert.Equal(t, float64(0), rate, "rate should be 0")
		} else {
			assert.Nil(t, err, "err should be nil")
			assert.Equal(t, tc.expectedRate, rate, "rate doesn't match the expected one")
		}
	}
}

func TestGetRate_EmptyRates(t *testing.T) {

	// Setup:
	rates := currencies.NewRates(time.Time{}, nil)

	// Verify:
	rate, err := rates.GetRate("USD", "EUR")

	assert.NotNil(t, err, "err shouldn't be nil")
	assert.Equal(t, float64(0), rate, "rate should be 0")
}
