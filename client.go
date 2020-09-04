package goeonet

import (
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var client HTTPClient

func init() {
	client = &http.Client{Timeout: 5 * time.Second}
}

func sendRequestToEonetApi(url string) ([]byte, error) {
	request, _ := http.NewRequest("GET", url, nil)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}