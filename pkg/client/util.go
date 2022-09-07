package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const storageAPIEndpoint = "api/v1/storage"

func (c FluentBitClient) GetMetricData() (*Response, error) {
	url := fmt.Sprintf("http://%s:%d/%s", c.FBHost, c.FBPort, storageAPIEndpoint)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error doing the request the HTTP request %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Received unexpected status code %d requesting %s", resp.StatusCode, storageAPIEndpoint)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response %v", err)
	}
	response := &Response{}
	jsonerr := json.Unmarshal([]byte(body), response)
	if jsonerr != nil {
		return nil, fmt.Errorf("Error unmarshalling FluentBit response %v", err)
	}
	return response, nil
}
