package goeonet

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	layoutISO         = "2006-01-02"
	baseEventsUrl     = "https://eonet.sci.gsfc.nasa.gov/api/v3/events"
	baseCategoriesUrl = "https://eonet.sci.gsfc.nasa.gov/api/v3/categories"
	baseLayersUrl     = "https://eonet.sci.gsfc.nasa.gov/api/v3/layers"
	baseSourcesUrl    = "https://eonet.sci.gsfc.nasa.gov/api/v3/sources"
)

type Category struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link,omitempty"`
	Description string `json:"description,omitempty"`
	Layers      string `json:"layers,omitempty"`
}

type Source struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Source string `json:"source"`
	Link   string `json:"link"`
}

type Sources struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Link        string   `json:"link"`
	Sources     []Source `json:"sources"`
}

type Coordinates struct {
	Coordinates [][]float64
}

func (c *Coordinates) UnmarshalJSON(data []byte) error {
	dataString := strings.Replace(string(data), " ", "", -1)
	dataString = strings.Replace(dataString, "],", "", -1)
	dataString = strings.Replace(dataString, "]", "", -1)
	dataString = strings.Replace(dataString, "[[", "", -1)
	coordinates := make([][]float64, 0)
	for _, coords := range strings.Split(dataString[1:], "[") {
		coordArr := make([]float64, 0)
		split := strings.Split(coords, ",")
		coord1, _ := strconv.ParseFloat(split[0], 64)
		coord2, _ := strconv.ParseFloat(split[1], 64)
		coordArr = append(coordArr, coord1)
		coordArr = append(coordArr, coord2)
		coordinates = append(coordinates, coordArr)
	}
	c.Coordinates = coordinates
	return nil
}

type EventSource struct {
	Id     string `json:"id"`
	Url    string `json:"url"`
}

type Geometry struct {
	MagnitudeValue float64     `json:"magnitudeValue"`
	MagnitudeUnit  string      `json:"magnitudeUnit"`
	Date           string      `json:"date"`
	Type           string      `json:"type"`
	Coordinates    Coordinates `json:"coordinates"`
}

type Event struct {
	Id          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Link        string        `json:"link"`
	Closed      string        `json:"closed"`
	Categories  []Category    `json:"categories"`
	Sources     []EventSource `json:"sources"`
	Geometrics  []Geometry    `json:"geometry"`
}

type Collection struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Link        string     `json:"link"`
	Events      []Event    `json:"events,omitempty"`
	Categories  []Category `json:"categories,omitempty"`
	Sources     []Source   `json:"sources,omitempty"`
}

type eventsQuery struct {
	source string
	status string
	limit  string
	days   string
	start  string
	end    string
	magID  string
	magMin string
	magMax string
	bbox   string
}

type categoriesQuery struct {
	category string
	source   string
	status   string
	limit    string
	days     string
}

var client = http.Client{Timeout: 5 * time.Second}

func GetRecentOpenEvents(limit string) (*Collection, error) {
	url := createEventsApiUrl(eventsQuery{limit: limit})

	collection, err := queryEventsApi(url.String())
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func GetEventsByDate(startDate, endDate string) (*Collection, error) {
	if !isValidDate(startDate) {
		return nil, errors.New("the starting date is invalid")
	}

	if endDate != "" && !isValidDate(endDate) {
		return nil, errors.New("the ending date is invalid")
	}

	url := createEventsApiUrl(eventsQuery{start: startDate, end: endDate})

	collection, err := queryEventsApi(url.String())
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func isValidDate(date string) bool {
	_, err := time.Parse(layoutISO, date)
	if err != nil {
		return false
	} else {
		return true
	}
}

func GetEventsBySourceID(sourceID string) (*Collection, error) {
	url := createEventsApiUrl(eventsQuery{source: sourceID})

	collection, err := queryEventsApi(url.String())
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func queryEventsApi(url string) (*Collection, error) {
	responseData, err := sendRequest(url)
	if err != nil {
		return nil, err
	}

	var collection Collection

	if err := json.Unmarshal(responseData, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func createEventsApiUrl(query eventsQuery) url.URL {
	u := url.URL {
		Scheme: "https",
		Host: "eonet.sci.gsfc.nasa.gov",
		Path: "/api/v3/events",
	}
	q := u.Query()
	q.Set("source", query.source)
	q.Set("status", query.status)
	q.Set("limit", query.limit)
	q.Set("days", query.days)
	if query.start != "" {
		q.Set("start", query.start)
		q.Set("end", query.end)
	}
	q.Set("magID", query.magID)
	q.Set("magMin", query.magMin)
	q.Set("magMax", query.magMax)
	q.Set("bbox", query.bbox)
	u.RawQuery = q.Encode()
	return u
}

func GetSources() (*Collection, error) {
	collection, err := querySourcesApi()
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func querySourcesApi() (*Collection, error) {
	responseData, err := sendRequest(baseSourcesUrl)
	if err != nil {
		return nil, err
	}

	var collection Collection

	if err := json.Unmarshal(responseData, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func GetCategories() (*Collection, error) {
	collection, err := queryCategoriesApi(baseCategoriesUrl)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func GetEventsByCategoryID(categoryID string) (*Collection, error) {
	url := createCategoriesApiUrl(categoriesQuery{category: categoryID})

	collection, err := queryCategoriesApi(url.String())
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func queryCategoriesApi(url string) (*Collection, error) {
	responseData, err := sendRequest(url)
	if err != nil {
		return nil, err
	}

	var collection Collection

	if err := json.Unmarshal(responseData, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func createCategoriesApiUrl(query categoriesQuery) url.URL {
	u := url.URL {
		Scheme: "https",
		Host: "eonet.sci.gsfc.nasa.gov",
		Path: "/api/v3/categories/" + query.category,
	}
	q := u.Query()
	q.Set("source", query.source)
	q.Set("status", query.status)
	q.Set("limit", query.limit)
	q.Set("days", query.days)
	u.RawQuery = q.Encode()
	return u
}

func queryLayersApi(query string) (*Collection, error) {
	responseData, err := sendRequest(baseLayersUrl + query)
	if err != nil {
		return nil, err
	}

	var collection Collection

	if err := json.Unmarshal(responseData, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

func sendRequest(url string) ([]byte, error) {
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
