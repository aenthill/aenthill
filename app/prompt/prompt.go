// Package prompt TODO.
package prompt

import (
	"github.com/manifoldco/promptui"
)

func ask(label string, defaultValue string) (string, error) {
	p := promptui.Prompt{
		Label:     label,
		Default:   defaultValue,
		AllowEdit: true,
	}

	r, err := p.Run()
	if err != nil {
		return "", err
	}

	return r, nil
}

/*func confirm(label string) (string, error) {
	p := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	r, err := p.Run()
	if err != nil {
		return "", err
	}

	return r, nil
}*/
