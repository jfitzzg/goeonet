package goeonet

import (
	"fmt"
	"net/url"
)

const (
	baseEventsUrl = "https://eonet.sci.gsfc.nasa.gov/api/v3/events"
	baseGeoJsonEventsUrl = "https://eonet.sci.gsfc.nasa.gov/api/v3/events/geojson"
)

// EventsQueryParameters is used for specifying the query parameters
// that can be passed to the GetEvents function. More information on
// the query parameters can be found at https://eonet.sci.gsfc.nasa.gov/docs/v3
type EventsQueryParameters struct {
	Source string
	Status string
	Limit  uint
	Days   uint
	Start  string
	End    string
	MagID  string
	MagMin string
	MagMax string
	Bbox   string
}

// GetEvents gets a list of events based on the query parameters passed in
func GetEvents(query EventsQueryParameters) ([]byte, error) {
	url := createEventsApiUrl(query)

	responseData, err := sendRequestToEonetApi(url.String())
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func createEventsApiUrl(query EventsQueryParameters) url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "eonet.sci.gsfc.nasa.gov",
		Path:   "/api/v3/events",
	}
	q := u.Query()
	if query.Source != "" {
		q.Set("source", query.Source)
	}
	if query.Status != "" {
		q.Set("status", query.Status)
	}
	if query.Limit > 0 {
		q.Set("limit", fmt.Sprint(query.Limit))
	}
	if query.Days > 0 {
		q.Set("days", fmt.Sprint(query.Days))
	}
	if query.Start != "" {
		q.Set("start", query.Start)
		if query.End != "" {
			q.Set("end", query.End)
		}
	}
	if query.MagID != "" {
		q.Set("magID", query.MagID)
	}
	if query.MagMin != "" {
		q.Set("magMin", query.MagMin)
	}
	if query.MagMax != "" {
		q.Set("magMax", query.MagMax)
	}
	if query.Bbox != "" {
		q.Set("bbox", query.Bbox)
	}
	u.RawQuery = q.Encode()
	return u
}

// GetGeoJsonEvents gets a list of events in GeoJSON format based on the query parameters passed in
func GetGeoJsonEvents(query EventsQueryParameters) ([]byte, error) {
	url := createGeoJsonEventsApiUrl(query)

	responseData, err := sendRequestToEonetApi(url.String())
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func createGeoJsonEventsApiUrl(query EventsQueryParameters) url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "eonet.sci.gsfc.nasa.gov",
		Path:   "/api/v3/events/geojson",
	}
	q := u.Query()
	if query.Source != "" {
		q.Set("source", query.Source)
	}
	if query.Status != "" {
		q.Set("status", query.Status)
	}
	if query.Limit > 0 {
		q.Set("limit", fmt.Sprint(query.Limit))
	}
	if query.Days > 0 {
		q.Set("days", fmt.Sprint(query.Days))
	}
	if query.Start != "" {
		q.Set("start", query.Start)
		if query.End != "" {
			q.Set("end", query.End)
		}
	}
	if query.MagID != "" {
		q.Set("magID", query.MagID)
	}
	if query.MagMin != "" {
		q.Set("magMin", query.MagMin)
	}
	if query.MagMax != "" {
		q.Set("magMax", query.MagMax)
	}
	if query.Bbox != "" {
		q.Set("bbox", query.Bbox)
	}
	u.RawQuery = q.Encode()
	return u
}
