// Package geo provides the ability to calculate duration between two locations.
package geo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

const (
	statusOk       = "OK"
	statusNotFound = "NOT_FOUND"

	geolocationURL = "https://www.googleapis.com/geolocation/v1/geolocate?key="
)

var (
	// ErrUnavailable is returned in cases where the underlying Google Maps API
	// did not return an error, but an unexpected response was received.
	ErrUnavailable = errors.New("duration unavailable")

	// ErrBadLocation is returned in cases where either the 'from' or 'to' address
	// could not be found.
	ErrBadLocation = errors.New("failed to find one of the provided locations")

	geolocationBody = bytes.NewBuffer([]byte(`{"considerIp": "true"}`))

	defaultAvoid = maps.AvoidTolls
)

// TravelMode dictates the type of travel when determining the duration.
type TravelMode maps.Mode

var (
	// Drive indicates a driving TravelMode.
	Drive = TravelMode(maps.TravelModeDriving)
	// Walk indicates a walking TravelMode.
	Walk = TravelMode(maps.TravelModeWalking)
	// Bike indicates a Biking TravelMode.
	Bike = TravelMode(maps.TravelModeBicycling)
	// Transit indicates public transit as a TravelMode.
	Transit = TravelMode(maps.TravelModeTransit)

	travelModeStrings = map[TravelMode]string{
		Drive:   "Drive",
		Walk:    "Walk",
		Bike:    "Bike",
		Transit: "Transit",
	}
)

// String returns a user friendly string representation of a TravelMode.
func (t TravelMode) String() string {
	return travelModeStrings[t]
}

// Router provides the ability to calculate travel duration between Routes.
type Router struct {
	apiKey string

	client Communicator
}

// NewRouter initializes and returns a Router with a Google Maps API key.
func NewRouter(apiKey string) (*Router, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Router{
		apiKey: apiKey,
		client: c,
	}, nil
}

// Duration returns the time it will take to travel between
// the From and To address.
func (r Router) Duration(from, to string, tm TravelMode) (*time.Duration, error) {
	req := maps.DistanceMatrixRequest{
		Origins:      []string{from},
		Destinations: []string{to},
		Mode:         maps.Mode(tm),
		Avoid:        defaultAvoid,
	}

	res, err := r.client.DistanceMatrix(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	for _, row := range res.Rows {
		for _, el := range row.Elements {
			switch el.Status {
			case statusNotFound:
				return nil, ErrBadLocation
			case statusOk:
				return &el.Duration, nil
			}
		}
	}

	return nil, ErrUnavailable
}

// CurrentLocation attempts to use Geolocation to return the Lat/Long of the system device
// based on it's IP Address.
func (r Router) CurrentLocation() (float64, float64, error) {
	resp, err := http.Post(geolocationURL+r.apiKey, "application/json", geolocationBody)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var loc struct {
		LatLng *maps.LatLng `json:"location"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&loc); err != nil {
		return 0, 0, err
	}

	return loc.LatLng.Lat, loc.LatLng.Lng, nil
}
