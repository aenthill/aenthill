package context

import (
	"fmt"
	"os"
	"testing"
)

// nolint: gocyclo
func TestNew(t *testing.T) {
	t.Run("calling New as if the application has been launched from host", func(t *testing.T) {
		if _, err := New(); err != nil {
			t.Errorf(`New should not have thrown an error: got "%s"`, err.Error())
		}
	})
	t.Run(fmt.Sprintf(`calling New as if the application has been launched from a container without "%s" environnment variable set`, IDEnvVar), func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo"}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be set`, IDEnvVar)
		}
		env = map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: ""}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be empty`, IDEnvVar)
		}
		unsetenv(t, env)
	})
	t.Run(fmt.Sprintf(`calling New as if the application has been launched from a container without "%s" environnment variable set`, FromContainerIDEnvVar), func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO"}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be set`, FromContainerIDEnvVar)
		}
		env = map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: ""}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be empty`, FromContainerIDEnvVar)
		}
		unsetenv(t, env)
	})
	t.Run(fmt.Sprintf(`calling New as if the application has been launched from a container without "%s" environnment variable set`, HostnameEnvVar), func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR"}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be set`, HostnameEnvVar)
		}
		env = map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: ""}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be empty`, HostnameEnvVar)
		}
		unsetenv(t, env)
	})
	t.Run(fmt.Sprintf(`calling New as if the application has been launched from a container without "%s" environnment variable set`, HostProjectDirEnvVar), func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ"}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be set`, HostProjectDirEnvVar)
		}
		env = map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ", HostProjectDirEnvVar: ""}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be empty`, HostProjectDirEnvVar)
		}
		unsetenv(t, env)
	})
	t.Run(fmt.Sprintf(`calling New as if the application has been launched from a container without "%s" environnment variable set`, ContainerProjectDirEnvVar), func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ", HostProjectDirEnvVar: "/foo"}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be set`, ContainerProjectDirEnvVar)
		}
		env = map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ", HostProjectDirEnvVar: "/foo", ContainerProjectDirEnvVar: ""}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be empty`, ContainerProjectDirEnvVar)
		}
		unsetenv(t, env)
	})
	t.Run(fmt.Sprintf(`calling New as if the application has been launched from a container without "%s" environnment variable set`, LogLevelEnvVar), func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ", HostProjectDirEnvVar: "/foo", ContainerProjectDirEnvVar: "/bar"}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be set`, LogLevelEnvVar)
		}
		env = map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ", HostProjectDirEnvVar: "/foo", ContainerProjectDirEnvVar: "/bar", LogLevelEnvVar: ""}
		setenv(t, env)
		if _, err := New(); err == nil {
			t.Errorf(`New should have thrown an error as "%s" should not be empty`, LogLevelEnvVar)
		}
		unsetenv(t, env)
	})
	t.Run("calling New as if the application has been launched from a container", func(t *testing.T) {
		env := map[string]string{ImageEnvVar: "aent/foo", IDEnvVar: "FOO", FromContainerIDEnvVar: "BAR", HostnameEnvVar: "BAZ", HostProjectDirEnvVar: "/foo", ContainerProjectDirEnvVar: "/bar", LogLevelEnvVar: "ERROR"}
		setenv(t, env)
		if _, err := New(); err != nil {
			t.Errorf(`New should not have thrown an error: got "%s"`, err.Error())
		}
		unsetenv(t, env)
	})
}

func setenv(t *testing.T, env map[string]string) {
	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf(`An unexpected error occurred while setting the environment variable "%s" with value "%s": got "%s"`, key, value, err.Error())
		}
	}
}

func unsetenv(t *testing.T, env map[string]string) {
	for key := range env {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf(`An unexpected error occurred while unsetting the environment variable "%s": got "%s"`, key, err.Error())
		}
	}
}

func TestIsContainer(t *testing.T) {
	ctx, err := New()
	if err != nil {
		t.Fatalf(`An unexpected error occurred while creating the context: got "%s"`, err.Error())
	}
	if ctx.IsContainer() {
		t.Error("IsContainer should have returned false as we are not in an aent")
	}
}
