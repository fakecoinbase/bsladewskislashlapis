package input

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/bsladewski/lapis/stream"
)

// A coinbaseStream is used to interact with the coinbase API.
type coinbaseStream struct {
	client http.Client
}

// NewCoinbaseStream retrieves a stream that can be used to retrieve values from
// the coinbase API.
func NewCoinbaseStream() stream.Stream {

	return &coinbaseStream{
		client: http.Client{Timeout: 15 * time.Second},
	}

}

// A coinbaseMockStream mocks interactions with the coinbase API.
type coinbaseMockStream struct {
	index      int
	spotPrices []float64
}

// NewCoinbaseMockStream retrieves a client that can be used to mock
// interactions with the coinbaes API.
func NewCoinbaseMockStream(mockDataReader io.Reader) (stream.Stream, error) {

	spotPrices, err := parseHistoricalData(mockDataReader)
	if err != nil {
		return nil, fmt.Errorf("parse mock data file, err: %v", err)
	}

	// return the mock client starting at index zero of historical price data
	return &coinbaseMockStream{
		index:      0,
		spotPrices: spotPrices,
	}, nil

}

// spotPriceResponse is used to read the exchange rate returned by the get spot
// price request to the coinbase API.
type spotPriceResponse struct {
	Amount string `json:"amount"`
}

func (d *coinbaseStream) Next() (float64, error) {

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

func (d *coinbaseStream) Close() {}

func (m *coinbaseMockStream) Next() (float64, error) {

	// retrieve the next item of mock data and increment current index into mock
	// data
	if m.index < len(m.spotPrices) {
		spotPrice := m.spotPrices[m.index]
		m.index++
		return spotPrice, nil
	}

	// return an error indicating that we have reached the end of mock data
	return 0.0, stream.ErrEndOfStream

}

func (m *coinbaseMockStream) Close() {
	m.spotPrices = nil
}

// GetHistoricalData retrieves historical hourly coinbase data for bitcoin
// prices.
func GetHistoricalData() ([]float64, error) {

	// create a new GET request for historical coinbase spot prices
	req, err := http.NewRequest(
		"GET",
		"http://www.cryptodatadownload.com/cdd/Coinbase_BTCUSD_1h.csv",
		nil)
	if err != nil {
		return nil, fmt.Errorf("get historical data, err: %v", err)
	}

	// create client for executing requests
	client := http.Client{Timeout: 30 * time.Second}

	// execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	spotPrices, err := parseHistoricalData(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse historical data, err: %v", err)
	}

	return spotPrices, nil

}

func parseHistoricalData(r io.Reader) ([]float64, error) {

	// define type for storing close prices associated with a specific timestamp
	type priceDatum struct {
		Timestamp time.Time
		Price     float64
	}

	// open csv reader
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	priceData := []priceDatum{}

	// loop until we are finished reading the csv data
	for {

		// read the next row
		record, err := reader.Read()

		if err == io.EOF {
			// if we are at the end of the file, the file has been successfully
			// read
			break
		} else if err != nil {
			// if any error other than an EOF is encountered, return the error
			// as we were unable to parse the csv file
			return nil, fmt.Errorf("reading mock data, err: %v", err)
		}

		// if we are not looking at a row of data, skip the row
		if record[1] != "BTCUSD" {
			continue
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

	return spotPrices, nil

}
