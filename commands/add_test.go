package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/manifest"
)

func TestNoImagesToAddError(t *testing.T) {
	err := &noImagesToAddError{}
	expected := fmt.Sprintf(noImagesToAddErrorMessage, RootCmd.Use, AddCmd.Use)

	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}

func TestAddCmd(t *testing.T) {
	// case 1: the manifest does not exist.
	if err := AddCmd.RunE(nil, nil); err == nil {
		t.Error("add command should have thrown an error because there is no manifest file")
	}

	// case 2: no args.
	copyManifest("aenthill-broken.json")
	if err := AddCmd.RunE(nil, nil); err == nil {
		t.Error("add command should have thrown an error because there is no args")
	}

	// case 3: the manifest is broken.
	if err := AddCmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
		t.Errorf("add command should have thrown an error because the manifest %s is broken", "aenthill-broken.json")
	}

	// case 4: the image does not exist.
	setFlags()
	copyManifest("aenthill-no-image.json")
	if err := AddCmd.RunE(nil, []string{"aent/foo"}); err == nil {
		t.Errorf("add command should have thrown an error because the image %s does not exist", "aent/foo")
	}

	// case 5: so far so good!
	if err := AddCmd.RunE(nil, []string{"aenthill/cassandra"}); err != nil {
		t.Error("add command should have worked")
	}

	// case 6: the image already exist.
	if err := AddCmd.RunE(nil, []string{"aenthill/cassandra"}); err != nil {
		t.Error("add command should have worked")
	}
	os.Remove(manifest.DefaultManifestFileName)
	unsetFlags()
}
