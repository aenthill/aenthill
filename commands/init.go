package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aenthill/manifest"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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

type manifestFileAlreadyExistingError struct{}

const manifestFileAlreadyExistingErrorMessage = "manifest %s already exists"

func (e *manifestFileAlreadyExistingError) Error() string {
	return fmt.Sprintf(manifestFileAlreadyExistingErrorMessage, manifest.DefaultManifestFileName)
}

/*
InitCmd creates the manifest in current directory.

Usage:

 aenthill init [flags]
*/
var InitCmd = &cobra.Command{
	Use:           "init",
	Short:         fmt.Sprintf("Create the manifest %s in current directory", manifest.DefaultManifestFileName),
	Long:          fmt.Sprintf("Create the manifest %s in current directory", manifest.DefaultManifestFileName),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(manifest.DefaultManifestFileName); err == nil {
			return &manifestFileAlreadyExistingError{}
		}

		// asking for project name.
		defaultName := ""
		wd, err := os.Getwd()
		if err == nil {
			defaultName = filepath.Base(wd)
		}

		name, err := ask("Project name", defaultName)
		if err != nil {
			return err
		}

		// asking for description.
		description, err := ask("Description", "")
		if err != nil {
			return err
		}

		return manifest.Flush(manifest.DefaultManifestFileName, &manifest.Manifest{Name: name, Description: description})
	},
}

func init() {
	RootCmd.AddCommand(InitCmd)
}
