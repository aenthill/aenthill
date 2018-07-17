// Package manifest eases the manipulation of an Aenthill manifest.
package manifest

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"regexp"

	"github.com/aenthill/aenthill/errors"
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
			Aents: make(map[string]*Aent),
		},
	}
}

type (
	data struct {
		Aents map[string]*Aent `json:"aents,omitempty"`
	}

	// Aent represents an entry from aents list.
	Aent struct {
		Image        string            `json:"image"`
		Metadata     map[string]string `json:"metadata,omitempty"`
		Events       []string          `json:"events,omitempty"`
		Dependencies map[string]string `json:"dependencies,omitempty"`
	}
)

var isAlpha = regexp.MustCompile(`^[A-Z0-9_]+$`).MatchString

// SetPath sets the path of the manifest file.
func (m *Manifest) SetPath(path string) {
	m.path = path
}

// Path returns the path of the manifest file.
func (m *Manifest) Path() string {
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
		return errors.Wrap("manifest", err)
	}
	return errors.Wrap("manifest", afero.WriteFile(m.fs, m.path, out, 0644))
}

// Parse simply parses the manifest file and populates the manifest data.
// Make sure your Manifest instance has a path to an existing file before
// using this function.
func (m *Manifest) Parse() error {
	data, err := afero.ReadFile(m.fs, m.path)
	if err != nil {
		return errors.Wrap("manifest", err)
	}
	return errors.Wrap("manifest", json.Unmarshal(data, &m.data))
}

// AddAent adds an aent in the manifest.
// Returns the generated key.
func (m *Manifest) AddAent(image string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	key := hex.EncodeToString(randBytes)
	m.data.Aents[key] = &Aent{Image: image}
	return key
}

// AddEvents adds events to an aent.
// If the key does not exist, throws an error.
func (m *Manifest) AddEvents(key string, events ...string) error {
	aent, ok := m.data.Aents[key]
	if !ok {
		return errors.Errorf("manifest", `aent identified by key "%s" does not exist`, key)
	}
	for _, event := range events {
		if !isAlpha(event) {
			return errors.Errorf("manifest", `"%s" is not a valid event name: only [A-Z0-9_] characters are authorized`, event)
		}
		if !m.isAentHandlingEvent(aent, event) {
			aent.Events = append(aent.Events, event)
		}
	}
	return nil
}

// AddMetadata adds metadata to an aent.
// If the key does not exist, throws an error.
func (m *Manifest) AddMetadata(key string, metadata map[string]string) error {
	aent, ok := m.data.Aents[key]
	if !ok {
		return errors.Errorf("manifest", `aent identified by key "%s" does not exist`, key)
	}
	if aent.Metadata == nil {
		aent.Metadata = make(map[string]string)
	}
	for k, value := range metadata {
		if !isAlpha(k) {
			return errors.Errorf("manifest", `"%s" is not a valid key for a metadata: only [A-Z0-9_] characters are authorized`, k)
		}
		aent.Metadata[k] = value
	}
	return nil
}

// Metadata returns the metadata of an aent.
// If the key does not exist, throws an error.
func (m *Manifest) Metadata(key string) (map[string]string, error) {
	aent, ok := m.data.Aents[key]
	if !ok {
		return nil, errors.Errorf("manifest", `aent identified by key "%s" does not exist`, key)
	}
	return aent.Metadata, nil
}

// AddDependency adds a dependency to an aent.
// Returns the dependency generated key.
// If the key does not exist or the dependency key does exist, throws an error.
func (m *Manifest) AddDependency(key, image, dependencyKey string) (string, error) {
	if !isAlpha(dependencyKey) {
		return "", errors.Errorf("manifest", `"%s" is not a valid key for a dependency: only [A-Z0-9_] characters are authorized`, dependencyKey)
	}
	aent, ok := m.data.Aents[key]
	if !ok {
		return "", errors.Errorf("manifest", `aent identified by key "%s" does not exist`, key)
	}
	if aent.Dependencies == nil {
		aent.Dependencies = make(map[string]string)
	}
	if _, ok := aent.Dependencies[dependencyKey]; ok {
		return "", errors.Errorf("manifest", `dependency identified by key "%s" does already exist for aent identified by key "%s"`, dependencyKey, key)
	}
	k := m.AddAent(image)
	aent.Dependencies[dependencyKey] = k
	return k, nil
}

// Dependencies returns the dependencies of an aent.
// If the key does not exist, throws an error.
func (m *Manifest) Dependencies(key string) (map[string]string, error) {
	aent, ok := m.data.Aents[key]
	if !ok {
		return nil, errors.Errorf("manifest", `aent identified by key "%s" does not exist`, key)
	}
	return aent.Dependencies, nil
}

// Aent returns an aent by its key.
// If the key does not exist, throws an error.
func (m *Manifest) Aent(key string) (*Aent, error) {
	aent, ok := m.data.Aents[key]
	if !ok {
		return nil, errors.Errorf("manifest", `aent identified by key "%s" does not exist`, key)
	}
	return aent, nil
}

// Aents returns a map of aents which handle the given event.
// If no event, returns all aents.
func (m *Manifest) Aents(event string) map[string]*Aent {
	if event == "" {
		return m.data.Aents
	}
	aents := make(map[string]*Aent)
	for key, aent := range m.data.Aents {
		if len(aent.Events) == 0 || m.isAentHandlingEvent(aent, event) {
			aents[key] = aent
		}
	}
	return aents
}

func (m *Manifest) isAentHandlingEvent(aent *Aent, event string) bool {
	for _, handled := range aent.Events {
		if handled == event {
			return true
		}
	}
	return false
}
