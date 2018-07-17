// Package errors contains useful function to create or wrap errors.
package errors

import "fmt"

// Error returns an error with the format "prefix: message".
func Error(prefix, message string) error {
	return fmt.Errorf("%s: %s", prefix, message)
}

// Errorf returns an error with the format: "prefix: formatted message".
func Errorf(prefix, format string, a ...interface{}) error {
	return fmt.Errorf("%s: %s", prefix, fmt.Sprintf(format, a...))
}

// Wrap returns an error with the format "prefix: error message".
func Wrap(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %s", prefix, err.Error())
}

// Wrapf returns an error with the format "prefix: formatted message: error message".
func Wrapf(prefix string, err error, format string, a ...interface{}) error {
	if err == nil {
		return err
	}
	return fmt.Errorf("%s: %s: %s", prefix, fmt.Sprintf(format, a...), err.Error())
}
