package manifest

import (
	"testing"

	"github.com/spf13/afero"
)

func TestNew(t *testing.T) {
	t.Run("calling New with a non-existing file", func(t *testing.T) {
		if m := New(DefaultManifestFileName, afero.NewMemMapFs()); m == nil {
			t.Error("Manifest should not be nil")
		}
	})
	t.Run("calling New with an existing file", func(t *testing.T) {
		if m := New("tests/aenthill.json", afero.NewOsFs()); m == nil {
			t.Error("Manifest should not be nil")
		}
	})
}

func TestSetPath(t *testing.T) {
	path := DefaultManifestFileName
	m := New("", afero.NewMemMapFs())
	m.SetPath(path)
	if m.path != path {
		t.Errorf(`SetPath set a wrong value: got "%s" want "%s"`, m.Path(), path)
	}
}

func TestPath(t *testing.T) {
	path := DefaultManifestFileName
	m := New(path, afero.NewMemMapFs())
	if m.Path() != path {
		t.Errorf(`Path returned a wrong value: got "%s" want "%s"`, m.Path(), path)
	}
}

func TestFlush(t *testing.T) {
	t.Run("calling Flush with an invalid file path", func(t *testing.T) {
		m := New("", afero.NewOsFs())
		if err := m.Flush(); err == nil {
			t.Error("Flush should have thrown an error as the given file path is not valid")
		}
	})
	t.Run("calling Flush with a valid file path", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Flush(); err != nil {
			t.Errorf(`Flush should not have thrown an error as the given file path is valid: got "%s"`, err.Error())
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("calling Parse with a non-existing file", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.Parse(); err == nil {
			t.Error("Parse should have thrown an error as file should not exist")
		}
	})
	t.Run("calling Parse with a broken file", func(t *testing.T) {
		m := New("../tests/aenthill-broken.json", afero.NewOsFs())
		if err := m.Parse(); err == nil {
			t.Error("Parse should have thrown an error as file should be broken")
		}
	})
	t.Run("calling Parse with a valid file", func(t *testing.T) {
		m := New("../tests/aenthill.json", afero.NewOsFs())
		if err := m.Parse(); err != nil {
			t.Errorf(`Parse should not have thrown an error as file should be valid: got "%s"`, err.Error())
		}
	})
}

func TestParseIfExist(t *testing.T) {
	t.Run("calling ParseIfExist with a non-existing file", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.ParseIfExist(); err != nil {
			t.Errorf(`ParseIfExist should not have thrown an error as file should not exist: got "%s"`, err.Error())
		}
	})
	t.Run("calling ParseIfExist with a broken file", func(t *testing.T) {
		m := New("../tests/aenthill-broken.json", afero.NewOsFs())
		if err := m.ParseIfExist(); err == nil {
			t.Error("ParseIfExist should have thrown an error as file should be broken")
		}
	})
	t.Run("calling ParseIfExist with a valid file", func(t *testing.T) {
		m := New("../tests/aenthill.json", afero.NewOsFs())
		if err := m.ParseIfExist(); err != nil {
			t.Errorf(`ParseIfExist should not have thrown an error as file should be valid: got "%s"`, err.Error())
		}
	})
}

func TestAddAent(t *testing.T) {
	m := New(DefaultManifestFileName, afero.NewMemMapFs())
	if key := m.AddAent("aent/foo"); key == "" {
		t.Error("AddAent should not have returned an empty key")
	}
}

func TestRemoveAent(t *testing.T) {
	m := New(DefaultManifestFileName, afero.NewMemMapFs())
	key := m.AddAent("aent/foo")
	if err := m.RemoveAent(key); err != nil {
		t.Fatalf(`RemoveAent should not have returned an error: got "%s"`, err.Error())
	}
	if _, err := m.Aent(key); err == nil {
		t.Error("RemoveAent should have removed the aent")
	}
}

func TestAddEvents(t *testing.T) {
	t.Run("calling AddEvents with a non-existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.AddEvents("foo", "FOO"); err == nil {
			t.Error("AddEvents should have thrown an error as given key should not exist")
		}
	})
	t.Run("calling AddEvents with an invalid event", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if err := m.AddEvents(key, "%FOO%"); err == nil {
			t.Error("AddEvents should have thrown an error as given event is not valid")
		}
	})
	t.Run("calling AddEvents with an existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if err := m.AddEvents(key, "FOO"); err != nil {
			t.Errorf(`AddEvents should not have thrown an error as given key should exist: got "%s"`, err.Error())
		}
	})
}

func TestAddMetadata(t *testing.T) {
	t.Run("calling AddMetadata with a non-existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.AddMetadata("foo", map[string]string{"FOO": "BAR"}); err == nil {
			t.Error("AddMetadata should have thrown an error as given key should not exist")
		}
	})
	t.Run("calling AddMetadata with an existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if err := m.AddMetadata(key, map[string]string{"FOO": "BAR"}); err != nil {
			t.Errorf(`AddMetadata should not have thrown an error as given key should exist: got "%s"`, err.Error())
		}
	})
}

func TestAddMetadataFlags(t *testing.T) {
	t.Run("calling AddMetadataFlags with empty flags", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.AddMetadataFromFlags("foo", nil); err != nil {
			t.Errorf(`AddMetadataFromFlags should not have thrown an error as flags are empty: got "%s"`, err.Error())
		}
	})
	t.Run("calling AddMetadataFlags with invalid flags", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if err := m.AddMetadataFromFlags("foo", []string{"FOO:BAR"}); err == nil {
			t.Error("AddMetadataFromFlags should have thrown an error as flags are invalid")
		}
	})
	t.Run("calling AddMetadataFlags with valid flags", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if err := m.AddMetadataFromFlags(key, []string{"FOO=BAR"}); err != nil {
			t.Errorf(`AddMetadataFromFlags should not have thrown an error as flags are valid: got "%s"`, err.Error())
		}
	})
}

func TestMetadata(t *testing.T) {
	t.Run("calling Metadata with a non-existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := m.Metadata("foo"); err == nil {
			t.Error("Metadata should have thrown an error as given key should not exist")
		}
	})
	t.Run("calling Metadata with an existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		metadata := make(map[string]string)
		metadata["FOO"] = "BAR"
		key := m.AddAent("aent/foo")
		if err := m.AddMetadata(key, metadata); err != nil {
			t.Errorf(`An unexpected error occurred while setting the metadata: got "%s"`, err.Error())
		}
		if _, err := m.Metadata(key); err != nil {
			t.Errorf(`Metadata should not have thrown an error as given key should exist: got "%s"`, err.Error())
		}
	})
}

func TestAddDependency(t *testing.T) {
	t.Run("calling AddDependency with a non-existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := m.AddDependency("foo", "aent/bar", "BAR"); err == nil {
			t.Error("AddDependency should have thrown an error as given key should not exist")
		}
	})
	t.Run("calling AddDependency with an existing dependency", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if _, err := m.AddDependency(key, "aent/bar", "BAR"); err != nil {
			t.Errorf(`An unexpected error occurred while setting the dependency: got "%s"`, err.Error())
		}
		if _, err := m.AddDependency(key, "aent/bar", "BAR"); err == nil {
			t.Errorf("AddDependency should have thrown an error as given dependency key already exist")
		}
	})
	t.Run("calling AddDependency with an existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if _, err := m.AddDependency(key, "aent/bar", "BAR"); err != nil {
			t.Errorf(`AddDependency should not have thrown an error: got "%s"`, err.Error())
		}
	})
}

func TestDependencies(t *testing.T) {
	t.Run("calling Dependencies with a non-existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := m.Dependencies("foo"); err == nil {
			t.Error("Dependencies should have thrown an error as given key should not exist")
		}
	})
	t.Run("calling Dependencies with an existing aent", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if _, err := m.Dependencies(key); err != nil {
			t.Errorf(`Dependencies should not have thrown an error as given key should exist: got "%s"`, err.Error())
		}
	})
}

func TestAent(t *testing.T) {
	t.Run("calling Aent with a non-existing key", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		if _, err := m.Aent("foo"); err == nil {
			t.Error("Aent should have thrown an error as given key should not exist")
		}
	})
	t.Run("calling Aent with an existing key", func(t *testing.T) {
		m := New(DefaultManifestFileName, afero.NewMemMapFs())
		key := m.AddAent("aent/foo")
		if _, err := m.Aent(key); err != nil {
			t.Error("Aent should nt have thrown an error as given key should exist")
		}
	})
}

func TestAents(t *testing.T) {
	t.Run("calling Aents without an event", func(t *testing.T) {
		m := New("../tests/aenthill.json", afero.NewOsFs())
		if err := m.Parse(); err != nil {
			t.Fatalf(`An unexpected error occurred while trying to parse the given manifest: got "%s"`, err.Error())
		}
		len := len(m.Aents(""))
		if len != 3 {
			t.Errorf(`Aents returned a map with a wrong length: got "%d" want "%d"`, len, 3)
		}
	})
	t.Run("calling Aents with a non-existing event", func(t *testing.T) {
		m := New("../tests/aenthill.json", afero.NewOsFs())
		if err := m.Parse(); err != nil {
			t.Fatalf(`An unexpected error occurred while trying to parse the given manifest: got "%s"`, err.Error())
		}
		len := len(m.Aents("BAZ"))
		if len != 2 {
			t.Errorf(`Aents returned a map with a wrong length: got "%d" want "%d"`, len, 2)
		}
	})
	t.Run("calling Aents with an existing event", func(t *testing.T) {
		m := New("../tests/aenthill.json", afero.NewOsFs())
		if err := m.Parse(); err != nil {
			t.Fatalf(`An unexpected error occurred while trying to parse the given manifest: got "%s"`, err.Error())
		}
		len := len(m.Aents("FOO"))
		if len != 3 {
			t.Errorf(`Aents returned a map with a wrong length: got "%d" want "%d"`, len, 3)
		}
	})
}
