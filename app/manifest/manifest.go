// Package manifest TODO.
package manifest

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	// Manifest TODO.
	Manifest struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Ants        []*Ant `yaml:"ants,omitempty"`
	}

	// Ant TODO.
	Ant struct {
		Image string `yaml:"image"`
	}
)

// DefaultManifestFileName TODO.
const DefaultManifestFileName = "aenthill.yml"

// Create TODO.
func Create(manifest *Manifest) error {
	f, err := os.Create(DefaultManifestFileName)
	if err != nil {
		return err
	}

	defer f.Close()

	out, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	if _, err := f.Write(out); err != nil {
		return err
	}

	return nil
}

// Parse TODO.
func Parse() (*Manifest, error) {
	m := &Manifest{}

	data, err := ioutil.ReadFile(DefaultManifestFileName)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}
