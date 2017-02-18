// Package geo provides the ability to calculate duration between two locations.
package geo

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

const (
	statusOk       string = "OK"
	statusNotFound string = "NOT_FOUND"
)

var (
	// ErrUnavailable is returned in cases where the underlying Google Maps API
	// did not return an error, but an unexpected response was received.
	ErrUnavailable = errors.New("duration unavailable")

	// ErrBadLocation is returned in cases where either the 'from' or 'to' address
	// could not be found.
	ErrBadLocation = errors.New("failed to find one of the provided locations")
)

// Router provides the ability to calculate travel duration between Routes.
type Router struct {
	client Communicator
}

// NewRouter initializes and returns a Router with a Google Maps API key.
func NewRouter(apiKey string) (*Router, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &Router{
		client: c,
	}, nil
}

// Duration returns the time it will take to travel between
// the From and To address.
func (r *Router) Duration(from, to string) (*time.Duration, error) {
	req := maps.DistanceMatrixRequest{
		Origins:      []string{from},
		Destinations: []string{to},
	}

	res, err := r.client.DistanceMatrix(context.Background(), &req)
	if err != nil {
		return nil, err
	}

	for _, row := range res.Rows {
		for _, el := range row.Elements {
			if el.Status == statusNotFound {
				return nil, ErrBadLocation
			}

			if el.Status == statusOk {
				return &el.Duration, nil
			}
		}
	}

	return nil, ErrUnavailable
}
