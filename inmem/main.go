package inmem

import (
	"encoding/json"
	"io/ioutil"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func throw(message string) error {
	return &errorString{message}
}

// InMem is Struct that represent InMem storage
type InMem struct {
	entries map[string]string
}

// New is method that creates an InMem Instance
func New() InMem {
	entries := make(map[string]string)
	return InMem{entries}
}

// Add is method to add new entry to db
func (im InMem) Add(key string, value string) (bool, error) {
	_, found := im.entries[key]

	if found {
		return false, throw("Key is already present")
	}

	im.entries[key] = value
	return true, nil
}

// Update is a method to update an entry
func (im InMem) Update(key string, value string) (bool, error) {
	_, found := im.entries[key]

	if found {
		im.entries[key] = value
		return true, nil
	}

	return false, throw("Key is not present")
}

// Remove is a method to remove a particular key
func (im InMem) Remove(key string) (bool, error) {
	_, found := im.entries[key]

	if !found {
		return false, throw("Key is not present")
	}

	delete(im.entries, key)
	return true, nil
}

// Get method gets the value of a key
func (im InMem) Get(key string) (string, bool) {
	value, found := im.entries[key]

	return value, found
}

// Persist writes inmem snapshot to a file location
func (im InMem) Persist(path string) error {
	b, err := json.Marshal(&im.entries)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0777)
}

// Entries returns all item in memory
func (im InMem) Entries() map[string]string {
	return im.entries
}

// Load snapshot from path to memory
func (im InMem) Load(path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, &im.entries)
}
