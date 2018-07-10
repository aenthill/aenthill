package errors

import "fmt"

func Error(prefix, message string) error {
	return fmt.Errorf("%s: %s", prefix, message)
}

func Errorf(prefix, format string, a ...interface{}) error {
	return fmt.Errorf("%s: %s", prefix, fmt.Sprintf(format, a...))
}

func Wrap(prefix string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %s", prefix, err.Error())
}

func Wrapf(prefix string, err error, format string, a ...interface{}) error {
	if err == nil {
		return err
	}
	return fmt.Errorf("%s: %s: %s", prefix, fmt.Sprintf(format, a...), err.Error())
}
