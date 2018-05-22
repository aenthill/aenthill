// Package manifest contains functions for manipulating an aenthill manifest.
package manifest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type (
	// Manifest gathers all data from an aenthill manifest.
	Manifest struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Aents       []*Aent `json:"aents,omitempty"`
	}

	// Aent represents an entry from aents list.
	Aent struct {
		Image string `json:"image"`
	}
)

// DefaultManifestFileName may be used in all functions with
// manifestFilePath as argument.
const DefaultManifestFileName = "aenthill.json"

// Flush populates the given file with manifest data.
// The data will be written as JSON.
func Flush(manifestFilePath string, manifest *Manifest) error {
	out, err := json.MarshalIndent(manifest, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(manifestFilePath, out, 0644)
}

// Parse simply parses the given file path and returns an instance of
// Manifest.
func Parse(manifestFilePath string) (*Manifest, error) {
	m := &Manifest{}

	data, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

type aentAlreadyInManifestError struct {
	image string
}

const aentAlreadyInManifestErrorMessage = "cannot add %s in given manifest as it does already exist"

func (e *aentAlreadyInManifestError) Error() string {
	return fmt.Sprintf(aentAlreadyInManifestErrorMessage, e.image)
}

// AddAent adds an image in given manifest.
// If the image does already exist, throws an error.
func AddAent(image string, manifest *Manifest) error {
	if getAentIndex(image, manifest) != -1 {
		return &aentAlreadyInManifestError{image}
	}

	manifest.Aents = append(manifest.Aents, &Aent{image})

	return nil
}

type aentNotInManifestError struct {
	image string
}

const aentNotInManifestErrorMessage = "cannot remove %s in given manifest as it does not exist"

func (e *aentNotInManifestError) Error() string {
	return fmt.Sprintf(aentNotInManifestErrorMessage, e.image)
}

// RemoveAent removes an image in given manifest.
// If the image does not exist, throws an error.
func RemoveAent(image string, manifest *Manifest) error {
	index := getAentIndex(image, manifest)
	if index == -1 {
		return &aentNotInManifestError{image}
	}

	manifest.Aents = append(manifest.Aents[:index], manifest.Aents[index+1:]...)

	return nil
}

// Exist returns true if the image exists in given manifest,
// false otherwise.
func Exist(image string, manifest *Manifest) bool {
	return getAentIndex(image, manifest) != -1
}

func getAentIndex(image string, manifest *Manifest) int {
	for i, aent := range manifest.Aents {
		if aent.Image == image {
			return i
		}
	}

	return -1
}
