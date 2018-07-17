package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewInstallJob(t *testing.T) {
	if j := NewInstallJob(nil, nil, nil, nil); j == nil {
		t.Error("NewInstallJob should not have returned an empty job")
	}
}

func TestInstallJobExecute(t *testing.T) {
	t.Run("calling Execute from install job with an invalid metadata", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		j := NewInstallJob([]string{"FOO:bar"}, nil, tests.MakeTestContext(t), m)
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given metadata is not valid")
		}
	})
	t.Run("calling Execute from install job with an invalid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		j := NewInstallJob(nil, []string{"%FOO%"}, tests.MakeTestContext(t), m)
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling Execute from install job with valid parameters", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		j := NewInstallJob([]string{"FOO=bar"}, []string{"FOO"}, tests.MakeTestContext(t), m)
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
		j = NewInstallJob(nil, []string{"FOO"}, tests.MakeTestContext(t), m)
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
	})
}
