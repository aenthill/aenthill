package docker

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"
)

func TestNew(t *testing.T) {
	if _, err := New(nil); err != nil {
		t.Errorf(`New should not have thrown an error: got "%s"`, err.Error())
	}
}

func TestRun(t *testing.T) {
	aent := &manifest.Aent{
		Image:        "aenthill/cassandra",
		Metadata:     map[string]string{"FOO": "BAR"},
		Dependencies: map[string]string{"FOO": "BAR"},
	}
	d, err := New(tests.MakeTestContext(t))
	if err != nil {
		t.Errorf(`An unexpected error occurred while creating a Docker instance: got "%s"`, err.Error())
	}
	if err := d.Run(aent, "", "FOO", ""); err != nil {
		t.Errorf(`Run should not have thrown an error as given image should exist: got "%s"`, err.Error())
	}
}

func TestReply(t *testing.T) {
	d, err := New(tests.MakeTestContext(t))
	if err != nil {
		t.Errorf(`An unexpected error occurred while creating a Docker instance: got "%s"`, err.Error())
	}
	if err := d.Reply("FOO", ""); err == nil {
		t.Error("Reply should have thrown an error as given container ID should not exist")
	}
}