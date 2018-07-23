// Package manifest eases the manipulation of an Aenthill manifest.
package manifest

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strings"

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

// SetPath sets the path of the manifest file.
func (m *Manifest) SetPath(path string) {
	m.path = path
}

// Path returns the path of the manifest file.
func (m *Manifest) Path() string {
	return m.path
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
	if err := m.exist(); err != nil {
		return err
	}
	return m.parse()
}

// ParseIfExist if exist parses the manifest file and populates the manifest data
// only if the manifest file does exist.
func (m *Manifest) ParseIfExist() error {
	if err := m.exist(); err != nil {
		return nil
	}
	return m.parse()
}

func (m *Manifest) parse() error {
	data, err := afero.ReadFile(m.fs, m.path)
	if err != nil {
		return errors.Wrap("manifest", err)
	}
	return errors.Wrap("manifest", json.Unmarshal(data, &m.data))
}

func (m *Manifest) exist() error {
	_, err := m.fs.Stat(m.path)
	return errors.Wrap("manifest", err)
}

var isAlpha = regexp.MustCompile(`^[A-Z0-9_]+$`).MatchString

// Validate checks if given string is only composed of [A-Z0-9_] characters.
func (m *Manifest) Validate(str, kind string) error {
	if isAlpha(str) {
		return nil
	}
	return errors.Errorf("manifest", `"%s" is not a valid %s: only [A-Z0-9_] characters are authorized`, str, kind)
}

// AddAent adds an aent in the manifest.
// Returns the generated ID.
func (m *Manifest) AddAent(image string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	ID := hex.EncodeToString(randBytes)
	m.data.Aents[ID] = &Aent{Image: image}
	return ID
}

// RemoveAent removes an aent from the manifest.
// If the ID does not exist, throws an error.
func (m *Manifest) RemoveAent(ID string) error {
	if _, err := m.Aent(ID); err != nil {
		return errors.Wrap("manifest", err)
	}
	delete(m.data.Aents, ID)
	return nil
}

// AddEvents adds events to an aent.
// If the ID does not exist, throws an error.
func (m *Manifest) AddEvents(ID string, events ...string) error {
	aent, err := m.Aent(ID)
	if err != nil {
		return err
	}
	for _, event := range events {
		if err := m.Validate(event, "event"); err != nil {
			return err
		}
		if !m.isAentHandlingEvent(aent, event) {
			aent.Events = append(aent.Events, event)
		}
	}
	return nil
}

// AddMetadata adds metadata to an aent.
// If the ID does not exist, throws an error.
func (m *Manifest) AddMetadata(ID string, metadata map[string]string) error {
	aent, err := m.Aent(ID)
	if err != nil {
		return err
	}
	if aent.Metadata == nil {
		aent.Metadata = make(map[string]string)
	}
	for k, value := range metadata {
		aent.Metadata[k] = value
	}
	return nil
}

// AddMetadataFromFlags adds metadata from flags to an aent.
// If the ID does not exist, throws an error.
func (m *Manifest) AddMetadataFromFlags(ID string, flags []string) error {
	if flags == nil {
		return nil
	}
	metadata := make(map[string]string)
	for _, data := range flags {
		parts := strings.Split(data, "=")
		if len(parts) != 2 {
			return errors.Errorf("manifest", `wrong metadata format: got "%s" want "key=value"`, data)
		}
		metadata[parts[0]] = parts[1]
	}
	return m.AddMetadata(ID, metadata)
}

// Metadata returns the metadata of an aent.
// If the ID does not exist, throws an error.
func (m *Manifest) Metadata(ID string) (map[string]string, error) {
	aent, err := m.Aent(ID)
	if err != nil {
		return nil, err
	}
	return aent.Metadata, nil
}

// AddDependency adds a dependency to an aent.
// Returns the dependency generated ID.
// If the ID does not exist or the dependency key does exist, throws an error.
func (m *Manifest) AddDependency(ID, image, key string) (string, error) {
	aent, err := m.Aent(ID)
	if err != nil {
		return "", err
	}
	if aent.Dependencies == nil {
		aent.Dependencies = make(map[string]string)
	}
	if _, ok := aent.Dependencies[key]; ok {
		return "", errors.Errorf("manifest", `dependency identified by key "%s" does already exist for aent identified by ID "%s"`, key, ID)
	}
	k := m.AddAent(image)
	aent.Dependencies[key] = k
	return k, nil
}

// Dependencies returns the dependencies of an aent.
// If the key does not exist, throws an error.
func (m *Manifest) Dependencies(key string) (map[string]string, error) {
	aent, err := m.Aent(key)
	if err != nil {
		return nil, err
	}
	return aent.Dependencies, nil
}

// Aent returns an aent by its ID.
// If the ID does not exist, throws an error.
func (m *Manifest) Aent(ID string) (*Aent, error) {
	aent, ok := m.data.Aents[ID]
	if !ok {
		return nil, errors.Errorf("manifest", `aent identified by ID "%s" does not exist`, ID)
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
	for ID, aent := range m.data.Aents {
		if len(aent.Events) == 0 || m.isAentHandlingEvent(aent, event) {
			aents[ID] = aent
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
