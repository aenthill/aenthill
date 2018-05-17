package manifest

import (
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	// case 1: creates the manifest.
	if err := Create(&Manifest{}); err != nil {
		t.Errorf("Manifest %s should have been created.", DefaultManifestFileName)
	}

	os.Remove(DefaultManifestFileName)
}

func TestParse(t *testing.T) {
	// case 1: the manifest does not exist.
	if _, err := Parse(); err == nil {
		t.Errorf("Manifest %s should not have been parsed as it should not exist.", DefaultManifestFileName)
	}

	// case 2: the manifest exists.
	Create(&Manifest{})
	if _, err := Parse(); err != nil {
		t.Errorf("Manifest %s should have been parsed as it should exist.", DefaultManifestFileName)
	}
}
