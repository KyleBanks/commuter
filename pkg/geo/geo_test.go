package geo

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type MockCommunicator struct {
	distanceFn func(context.Context, *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error)
}

func (m *MockCommunicator) DistanceMatrix(c context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
	return m.distanceFn(c, r)
}

func TestNewRouter(t *testing.T) {
	if _, err := NewRouter(""); err == nil {
		t.Fatal("Expected error for empty API key")
	}

	r, err := NewRouter("apiKey")
	if err != nil {
		t.Fatal(err)
	}

	if r.client == nil {
		t.Fatal("Unexpected nil client")
	}
}

func TestRouter_Duration(t *testing.T) {
	var mc MockCommunicator
	r := Router{
		client: &mc,
	}

	// Positive Case
	{
		from := "321 Main St"
		to := "123 Maple St"
		duration := time.Minute * 5
		mode := Drive
		mc.distanceFn = func(c context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
			if c == nil {
				t.Fatal("Unexpected nil Context")
			}

			if len(r.Origins) != 1 || r.Origins[0] != from {
				t.Fatalf("Unexpected Origins, expected=%v, got=%v", from, r.Origins)
			}
			if len(r.Destinations) != 1 || r.Destinations[0] != to {
				t.Fatalf("Unexpected Destinations, expected=%v, got=%v", to, r.Destinations)
			}
			if r.Mode != maps.Mode(mode) {
				t.Fatalf("Unexpected Mode, expected=%v, got=%v", mode, r.Mode)
			}
			if r.Avoid != defaultAvoid {
				t.Fatalf("Unexpected Avoid, expected=%v, got=%v", defaultAvoid, r.Avoid)
			}

			return &maps.DistanceMatrixResponse{
				Rows: []maps.DistanceMatrixElementsRow{
					{
						Elements: []*maps.DistanceMatrixElement{
							&maps.DistanceMatrixElement{
								Status:   statusOk,
								Duration: duration,
							},
						},
					},
				},
			}, nil
		}

		d, err := r.Duration(from, to, mode)
		if err != nil {
			t.Fatal(err)
		}

		if *d != duration {
			t.Fatalf("Unexpected duration returned, expected=%v, got=%v", duration, *d)
		}
	}

	// Error from Communicator
	{
		e := errors.New("test err")
		mc.distanceFn = func(c context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
			return nil, e
		}

		_, err := r.Duration("", "", Drive)
		if err != e {
			t.Fatalf("Unexpected error returned, expected=%v, got=%v", e, err)
		}
	}

	// Error Status
	{
		mc.distanceFn = func(c context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
			return &maps.DistanceMatrixResponse{
				Rows: []maps.DistanceMatrixElementsRow{
					{
						Elements: []*maps.DistanceMatrixElement{
							&maps.DistanceMatrixElement{
								Status: statusNotFound,
							},
						},
					},
				},
			}, nil
		}

		_, err := r.Duration("", "", Drive)
		if err != ErrBadLocation {
			t.Fatalf("Unexpected error returned, expected=%v, got=%v", ErrBadLocation, err)
		}
	}

	// No response
	{
		mc.distanceFn = func(c context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
			return &maps.DistanceMatrixResponse{}, nil
		}

		_, err := r.Duration("", "", Drive)
		if err != ErrUnavailable {
			t.Fatalf("Unexpected error returned, expected=%v, got=%v", ErrUnavailable, err)
		}
	}
}
