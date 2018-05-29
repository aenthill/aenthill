// Package manifest is a library easing the manipulation of an aenthill manifest.
package manifest

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/afero"
)

// DefaultManifestFileName may be used as path argument for
// NewManifest function.
const DefaultManifestFileName = "aenthill.json"

// Manifest is our working struct for manipulating an aenthill manifest.
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
		Image string `json:"image"`
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

type aentAlreadyInManifestError struct {
	image string
}

const aentAlreadyInManifestErrorMessage = "cannot add %s in given manifest as it does already exist"

func (e *aentAlreadyInManifestError) Error() string {
	return fmt.Sprintf(aentAlreadyInManifestErrorMessage, e.image)
}

// AddAent adds an image in the manifest.
// If the image does already exist, throws an error.
func (m *Manifest) AddAent(image string) error {
	if m.getAentIndex(image) != -1 {
		return &aentAlreadyInManifestError{image}
	}
	m.data.Aents = append(m.data.Aents, &Aent{image})

	return nil
}

type aentNotInManifestError struct {
	image string
}

const aentNotInManifestErrorMessage = "cannot remove %s in given manifest as it does not exist"

func (e *aentNotInManifestError) Error() string {
	return fmt.Sprintf(aentNotInManifestErrorMessage, e.image)
}

// RemoveAent removes an image from the manifest.
// If the image does not exist, throws an error.
func (m *Manifest) RemoveAent(image string) error {
	index := m.getAentIndex(image)
	if index == -1 {
		return &aentNotInManifestError{image}
	}
	m.data.Aents = append(m.data.Aents[:index], m.data.Aents[index+1:]...)

	return nil
}

// HasAent returns true if the image exists in the manifest,
// false otherwise.
func (m *Manifest) HasAent(image string) bool {
	return m.getAentIndex(image) != -1
}

// GetAents returns the array of aents of the manifest.
func (m *Manifest) GetAents() []*Aent {
	return m.data.Aents
}

func (m *Manifest) getAentIndex(image string) int {
	for i, aent := range m.data.Aents {
		if aent.Image == image {
			return i
		}
	}

	return -1
}
