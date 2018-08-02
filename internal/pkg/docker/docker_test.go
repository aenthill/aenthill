package docker

import (
	"testing"

	"github.com/aenthill/aenthill/test"
)

func TestNew(t *testing.T) {
	if _, err := New(nil); err != nil {
		t.Errorf(`New should not have thrown an error: got "%s"`, err.Error())
	}
}

func TestRun(t *testing.T) {
	d, err := New(test.Context(t))
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating a Docker instance: got "%s"`, err.Error())
	}
	if err := d.Run("aenthill/cassandra", "", "FOO", ""); err != nil {
		t.Errorf(`Run should not have thrown an error as given image should exist: got "%s"`, err.Error())
	}
}

func TestReply(t *testing.T) {
	d, err := New(test.Context(t))
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating a Docker instance: got "%s"`, err.Error())
	}
	if err := d.Reply("FOO", ""); err == nil {
		t.Error("Reply should have thrown an error as given container ID should not exist")
	}
}
