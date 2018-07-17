package jobs

import (
	"testing"

	"github.com/aenthill/aenthill/manifest"
	"github.com/aenthill/aenthill/tests"

	"github.com/spf13/afero"
)

func TestNewRegisterJob(t *testing.T) {
	if j := NewRegisterJob("", "", nil, nil, nil, nil); j == nil {
		t.Error("NewRegisterJob should not have returned an empty job")
	}
}

func TestRegisterJobExecute(t *testing.T) {
	t.Run("calling Execute from register job with no key in context", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		j := NewRegisterJob("aent/foo", "FOO", nil, nil, tests.MakeTestContext(t), m)
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as there is no key in context")
		}
	})
	t.Run("calling Execute from register job with an invalid metadata", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		ctx := tests.MakeTestContext(t)
		ctx.Key = key
		j := NewRegisterJob("aent/bar", "FOO", []string{"FOO:bar"}, nil, ctx, m)
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given metadata is not valid")
		}
	})
	t.Run("calling Execute from register job with an invalid event", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		ctx := tests.MakeTestContext(t)
		ctx.Key = key
		j := NewRegisterJob("aent/bar", "FOO", nil, []string{"%FOO%"}, ctx, m)
		if err := j.Execute(); err == nil {
			t.Error("Execute should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling Execute from register job with valid parameters", func(t *testing.T) {
		m := manifest.New(manifest.DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		ctx := tests.MakeTestContext(t)
		ctx.Key = key
		j := NewRegisterJob("aent/bar", "FOO", []string{"FOO=bar"}, []string{"FOO"}, ctx, m)
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
		j = NewRegisterJob("aent/baz", "BAR", nil, []string{"FOO"}, ctx, m)
		if err := j.Execute(); err != nil {
			t.Errorf(`Execute should not have thrown an error: got "%s"`, err.Error())
		}
	})
}
