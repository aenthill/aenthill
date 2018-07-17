package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/tests"
)

func TestNewReplyJob(t *testing.T) {
	t.Run("calling NewReplyJob an invalid event", func(t *testing.T) {
		if _, err := NewReplyJob("%FOO%", "", nil); err == nil {
			t.Error("NewReplyJob should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling NewReplyJob with a valid event", func(t *testing.T) {
		if _, err := NewReplyJob("FOO", "", nil); err != nil {
			t.Errorf(`NewReplyJob should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

func TestReplyJobExecute(t *testing.T) {
	j, err := NewReplyJob("FOO", "", tests.MakeTestContext(t))
	if err != nil {
		t.Errorf(`An unexpected error occurred while creating a reply job: got "%s"`, err.Error())
	}
	if err := j.Execute(); err == nil {
		t.Error("Execute should have thrown an error as given container ID should not exist")
	}
}
