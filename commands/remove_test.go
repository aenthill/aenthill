package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/aenthill/manifest"
)

func TestNoImagesToRemoveError(t *testing.T) {
	err := &noImagesToRemoveError{}
	expected := fmt.Sprintf(noImagesToRemoveErrorMessage, RootCmd.Use, RemoveCmd.Use)

	if err.Error() != expected {
		t.Errorf("error returned a wrong message: got %s want %s", err.Error(), expected)
	}
}

func TestRemoveCmd(t *testing.T) {
	// case 1: the manifest does not exist.
	if err := RemoveCmd.RunE(nil, nil); err == nil {
		t.Error("remove command should have thrown an error because there is no manifest file")
	}

	// case 2: no args.
	copyManifest("aenthill-broken.json")
	if err := RemoveCmd.RunE(nil, nil); err == nil {
		t.Error("remove command should have thrown an error because there is no args")
	}

	// case 3: the manifest is broken.
	if err := RemoveCmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
		t.Errorf("remove command should have thrown an error because the manifest %s is broken", "aenthill-broken.json")
	}

	// case 4: the image does not exist in manifest.
	copyManifest("aenthill-inexistant-image.json")
	if err := RemoveCmd.RunE(nil, []string{"aenthill/cassandra"}); err == nil {
		t.Errorf("remove command should have thrown an error because the image %s does not exist in %s", "aenthill/cassandra", "aenthill-no-image.json")
	}

	// case 5: the image does not exist.
	setFlags()
	if err := RemoveCmd.RunE(nil, []string{"aent/foo"}); err == nil {
		t.Errorf("remove command should have thrown an error because the image %s does not exist", "aent/foo")
	}

	// case 6: so far so good!
	copyManifest("aenthill.json")
	if err := RemoveCmd.RunE(nil, []string{"aenthill/cassandra"}); err != nil {
		t.Error("remove command should have worked")
	}
	os.Remove(manifest.DefaultManifestFileName)
	unsetFlags()

}
