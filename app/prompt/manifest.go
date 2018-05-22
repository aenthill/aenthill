package prompt

import (
	"os"
	"path/filepath"

	"github.com/aenthill/manifest"
)

// AskManifestValues retrieves user inputs in order to initialize a manifest.
func AskManifestValues() (*manifest.Manifest, error) {
	// asking for project name.
	defaultName := ""
	wd, err := os.Getwd()
	if err == nil {
		defaultName = filepath.Base(wd)
	}

	name, err := ask("Project name", defaultName)
	if err != nil {
		return nil, err
	}

	// asking for description.
	description, err := ask("Description", "")
	if err != nil {
		return nil, err
	}

	return &manifest.Manifest{Name: name, Description: description}, nil
}
