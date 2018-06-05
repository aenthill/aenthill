// Package manifest is a library easing the manipulation of an Aenthill manifest.
package manifest

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

// DefaultManifestFileName may be used as path argument for
// NewManifest function.
const DefaultManifestFileName = "aenthill.json"

// Manifest is our working struct for manipulating an Aenthill manifest.
type Manifest struct {
	path string
	fs   afero.Fs
	data *data
}

// New creates a Manifest instance with
// the given file path (may not exist) and file system.
func New(path string, fs afero.Fs) *Manifest {
	return &Manifest{
		path: path,
		fs:   fs,
		data: &data{
			Aents: make([]*Aent, 0),
		},
	}
}

type (
	data struct {
		Aents []*Aent `json:"aents,omitempty"`
	}

	// Aent represents an entry from aents list.
	Aent struct {
		Image         string   `json:"image"`
		HandledEvents []string `json:"handled_events,omitempty"`
	}
)

// GetPath returns the path of the manifest file.
func (m *Manifest) GetPath() string {
	return m.path
}

// Exist returns true if the manifest file exists.
func (m *Manifest) Exist() bool {
	_, err := m.fs.Stat(m.path)
	return err == nil
}

// Flush writes the manifest file and populates it with the manifest data.
// The data will be written as JSON.
func (m *Manifest) Flush() error {
	out, err := json.MarshalIndent(m.data, "", "\t")
	if err != nil {
		return err
	}

	return afero.WriteFile(m.fs, m.path, out, 0644)
}

// Parse simply parses the manifest file and populates the manifest data.
// Make sure your Manifest instance has a path to an existing file before
// using this function.
func (m *Manifest) Parse() error {
	data, err := afero.ReadFile(m.fs, m.path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &m.data); err != nil {
		return err
	}

	return nil
}

type addExistingAentError struct {
	image string
}

const addExistingAentErrorMessage = "cannot add %s in given manifest as it does already exist"

func (e *addExistingAentError) Error() string {
	return fmt.Sprintf(addExistingAentErrorMessage, e.image)
}

// AddAent adds an image in the manifest.
// If the image does already exist, throws an error.
func (m *Manifest) AddAent(image string) error {
	if m.getAentIndex(image) != -1 {
		return &addExistingAentError{image}
	}
	m.data.Aents = append(m.data.Aents, &Aent{Image: image})

	return nil
}

type setHandledEventsToNonExistingAentError struct {
	image string
}

const setHandledEventsToNonExistingAentErrorMessage = "cannot set handled events to %s in given manifest as it does not exist"

func (e *setHandledEventsToNonExistingAentError) Error() string {
	return fmt.Sprintf(setHandledEventsToNonExistingAentErrorMessage, e.image)
}

// SetHandledEvents sets the handled events of an image (previous handled events will be deleted).
// If the image does not exist, throws an error.
func (m *Manifest) SetHandledEvents(image string, events ...string) error {
	index := m.getAentIndex(image)
	if index == -1 {
		return &setHandledEventsToNonExistingAentError{image}
	}
	m.data.Aents[index].HandledEvents = append([]string{}, events...)

	return nil
}

type removeNonExistingAentError struct {
	image string
}

const removeNonExistingAentErrorMessage = "cannot remove %s in given manifest as it does not exist"

func (e *removeNonExistingAentError) Error() string {
	return fmt.Sprintf(removeNonExistingAentErrorMessage, e.image)
}

// RemoveAent removes an image from the manifest.
// If the image does not exist, throws an error.
func (m *Manifest) RemoveAent(image string) error {
	index := m.getAentIndex(image)
	if index == -1 {
		return &removeNonExistingAentError{image}
	}
	m.data.Aents = append(m.data.Aents[:index], m.data.Aents[index+1:]...)

	return nil
}

// HasAent returns true if the image exists in the manifest,
// false otherwise.
func (m *Manifest) HasAent(image string) bool {
	return m.getAentIndex(image) != -1
}

// GetAents returns the array of aents of the manifest which handle the event.
// If no event, returns all aents.
func (m *Manifest) GetAents(event string) []*Aent {
	if event == "" {
		return m.data.Aents
	}

	var aents []*Aent
	for _, aent := range m.data.Aents {
		if len(aent.HandledEvents) == 0 || m.isAentHandlingEvent(aent, event) {
			aents = append(aents, aent)
		}
	}

	return aents
}

func (m *Manifest) isAentHandlingEvent(aent *Aent, event string) bool {
	for _, handledEvent := range aent.HandledEvents {
		if handledEvent == event {
			return true
		}
	}

	return false
}

func (m *Manifest) getAentIndex(image string) int {
	for i, aent := range m.data.Aents {
		if aent.Image == image {
			return i
		}
	}

	return -1
}
