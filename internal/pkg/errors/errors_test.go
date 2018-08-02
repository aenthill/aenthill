package errors

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := Error("foo", "bar")
	expected := "foo: bar"
	if err.Error() != expected {
		t.Errorf(`error returned a wrong message: got "%s" want "%s"`, err.Error(), expected)
	}
}

func TestErrorf(t *testing.T) {
	err := Errorf("foo", "bar %s", "baz")
	expected := fmt.Sprintf("foo: bar %s", "baz")
	if err.Error() != expected {
		t.Errorf(`error returned a wrong message: got "%s" want "%s"`, err.Error(), expected)
	}
}

func TestWrap(t *testing.T) {
	t.Run("calling Wrap with a non-existing error", func(t *testing.T) {
		err := Wrap("foo", nil)
		if err != nil {
			t.Errorf(`Wrap should have returned nil: got "%s"`, err.Error())
		}
	})
	t.Run("calling Wrap with an existing error", func(t *testing.T) {
		err := Wrap("foo", Error("bar", "baz"))
		expected := "foo: bar: baz"
		if err.Error() != expected {
			t.Errorf(`error returned a wrong message: got "%s" want "%s"`, err.Error(), expected)
		}
	})
}

func TestWrapf(t *testing.T) {
	t.Run("calling Wrapf with a non-existing error", func(t *testing.T) {
		err := Wrapf("foo", nil, "foo bar %s", "baz")
		if err != nil {
			t.Errorf(`Wrapf should have returned nil: got "%s"`, err.Error())
		}
	})
	t.Run("calling Wrapf with an existing error", func(t *testing.T) {
		err := Wrapf("foo", Error("bar", "baz"), "foo bar %s", "baz")
		expected := fmt.Sprintf("foo: foo bar baz: bar: baz")
		if err.Error() != expected {
			t.Errorf(`error returned a wrong message: got "%s" want "%s"`, err.Error(), expected)
		}
	})
}
