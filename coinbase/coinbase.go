// Package coinbase provides a wrapper to the coinbase API.
package coinbase

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
	// ErrEndOfMockData indicates that the end of the historical data used to
	// mock the coinbase API has been reached.
	ErrEndOfMockData = errors.New("end of mock data")
)

// Client provides functions for interacting with the coinbase API.
type Client interface {
	// GetSpotPrice retrieves the exchange rate between one BTC and USD.
	GetSpotPrice() (float64, error)
}

// A defaultClient is used to interact with the coinbase API.
type defaultClient struct {
	client http.Client
}

// NewClient retrieves a client that can be used to interact with the coinbase
// API.
func NewClient() Client {

	return &defaultClient{
		client: http.Client{Timeout: 15 * time.Second},
	}

}

// A mockClient mocks interactions with the coinbase API.
type mockClient struct {
	index      int
	spotPrices []float64
}

// NewMockClient retrieves a client that can be used to mock interactions with
// the coinbaes API.
func NewMockClient() (Client, error) {

	// open file containing mock data
	mockDataFile, err := os.Open("mock_data.csv")
	if err != nil {
		return nil, fmt.Errorf("open mock data file, err: %v", err)
	}
	defer mockDataFile.Close()

	// read mock data as comma separated values
	reader := csv.NewReader(mockDataFile)

	// define type for storing close prices associated with a specific timestamp
	type priceDatum struct {
		Timestamp time.Time
		Price     float64
	}

	priceData := []priceDatum{}

	// headers indicates that the record we are currently looking at is a row of
	// headers and not data
	headers := true

	// loop until we are finished reading the csv data
	for {

		// read the next row
		record, err := reader.Read()

		// if we are looking at headers, skip the current row and mark that
		// headers have already been read
		if headers {
			headers = false
			continue
		}

		if err == io.EOF {
			// if we are at the end of the file, the file has been successfully
			// read
			break
		} else if err != nil {
			// if any error other than an EOF is encountered, return the error
			// as we were unable to parse the csv file
			return nil, fmt.Errorf("reading mock data, err: %v", err)
		}

		// parse the timestamp from the csv data; as the data starts with the
		// newest record we will need to sort by timestamp ascending when the
		// file has been read
		timestamp, err := time.Parse("2006-01-02 03-PM", record[0])
		if err != nil {
			return nil, fmt.Errorf("parsing timestamp, err: %v", err)
		}

		// parse the close price from csv data; the close price will represent
		// the spot price for this index
		price, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("parsing price, err: %v", err)
		}

		// add the timestamp and close price to price data
		priceData = append(priceData, priceDatum{timestamp, price})

	}

	// sort price data by timestamp ascending
	sort.Slice(priceData, func(i, j int) bool {
		return priceData[i].Timestamp.Before(priceData[j].Timestamp)
	})

	// build ordered list of price data
	spotPrices := []float64{}
	for _, priceDatum := range priceData {
		spotPrices = append(spotPrices, priceDatum.Price)
	}

	// return the mock client starting at index zero of historical price data
	return &mockClient{
		index:      0,
		spotPrices: spotPrices,
	}, nil

}

// spotPriceResponse is used to read the exchange rate returned by the get spot
// price request to the coinbase API.
type spotPriceResponse struct {
	Amount string `json:"amount"`
}

func (d *defaultClient) GetSpotPrice() (float64, error) {

	// create a new GET request to coinbase spot price endpoint
	req, err := http.NewRequest(
		"GET",
		"https://api.coinbase.com/v2/prices/spot?currency=USD",
		nil)
	if err != nil {
		return 0.0, err
	}

	// execute the request
	resp, err := d.client.Do(req)
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

func (m *mockClient) GetSpotPrice() (float64, error) {

	// retrieve the next item of mock data and increment current index into mock
	// data
	if m.index < len(m.spotPrices) {
		spotPrice := m.spotPrices[m.index]
		m.index++
		return spotPrice, nil
	}

	// return an error indicating that we have reached the end of mock data
	return 0.0, ErrEndOfMockData

}
