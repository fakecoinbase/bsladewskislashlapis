// Package coinbase provides a wrapper to the coinbase API.
package coinbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// client is used to execute requests to the coinbase API.
var client = http.Client{
	Timeout: 30 * time.Second,
}

// spotPriceResponse is used to read the exchange rate returned by the get spot
// price request to the coinbase API.
type spotPriceResponse struct {
	Amount string `json:"amount"`
}

// GetSpotPrice retrieves the exchange rate between one BTC and USD.
func GetSpotPrice() (float64, error) {

	// create a new GET request to coinbase spot price endpoint
	req, err := http.NewRequest(
		"GET",
		"https://api.coinbase.com/v2/prices/spot?currency=USD",
		nil)
	if err != nil {
		return 0.0, err
	}

	// execute the request
	resp, err := client.Do(req)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		// if request was not successfull return error containing status code
		// and response body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0.0, err
		}

		return 0.0, fmt.Errorf("%d - %x", resp.StatusCode, respBody)

	}

	// parse the exchange rate from the response
	var respData struct {
		Data spotPriceResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return 0.0, err
	}

	rate, err := strconv.ParseFloat(respData.Data.Amount, 64)
	if err != nil {
		return 0.0, err
	}

	return rate, nil

}
