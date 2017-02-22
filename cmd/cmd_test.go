package cmd

import (
	"fmt"
	"time"
)

// mock Indicator

type mockIndicator struct {
	out []string
}

func (m *mockIndicator) Indicate(msg string, args ...interface{}) {
	m.out = append(m.out, fmt.Sprintf(msg, args...))
}

// mock Durationer

type mockDurationer struct {
	durationFn func(string, string) (*time.Duration, error)
}

func (m *mockDurationer) Duration(from, to string) (*time.Duration, error) {
	return m.durationFn(from, to)
}

// mock StorageProvider

type mockStorageProvider struct {
	loadFn func(interface{}) error
	saveFn func(interface{}) error
}

func (m *mockStorageProvider) Load(i interface{}) error {
	return m.loadFn(i)
}

func (m *mockStorageProvider) Save(i interface{}) error {
	return m.saveFn(i)
}

// mock Locator

type mockLocator struct {
	locateFn func() (float64, float64, error)
}

func (m *mockLocator) CurrentLocation() (float64, float64, error) {
	return m.locateFn()
}
