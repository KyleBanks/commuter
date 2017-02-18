// Package storage provides the ability to persist and retrieve structs.
package storage

import (
	"encoding/json"
	"io/ioutil"
)

// Provider defines a type that can be used for storage.
type Provider interface {
	Load(interface{}) error
	Save(interface{}) error
}

// FileStore represents a simple on-disk file storage system.
type FileStore struct {
	Filename string
}

// NewFileStore returns an initialized FileStore.
func NewFileStore(f string) *FileStore {
	return &FileStore{
		Filename: f,
	}
}

// Load attempts to read and decode the storage file into the
// value provided.
func (f FileStore) Load(v interface{}) error {
	data, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		return err
	}

	json.Unmarshal(data, v)
	return nil
}

// Save attempts to write the value provided out to the storage file.
func (f FileStore) Save(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.Filename, data, 0644)
}
