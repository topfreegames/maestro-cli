// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package init

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/extensions"
)

// initCmd represents the init maestro-cli command
var Cmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize maestro-cli",
	Long:  `Creates the directory ~/.maestro and appends a default config file`,
	Args:  validateArgs,
	Run:   run,
}

func validateArgs(_ *cobra.Command, args []string) error {

	if len(args) == 0 {
		return errors.New("missing arg with maestro-cli context name")
	} else if len(args) < 2 {
		return errors.New("missing arg with maestro server URL")
	}

	_, err := url.ParseRequestURI(args[1])
	if err != nil {
		return errors.New("bad maestro server URl")
	}

	return nil
}

func run(_ *cobra.Command, args []string) {
	context := args[0]
	serverURL := args[1]
	config := extensions.NewConfig(serverURL)
	err := config.Write(extensions.NewFileSystem(), context)
	if err != nil {
		fmt.Println("Error writing config file: ", err)
		os.Exit(1)
	}

	fmt.Println("Configuration created")
}
