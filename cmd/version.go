// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"

	"github.com/topfreegames/maestro-cli/metadata"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Maestro-cli version",
	Long:  `Maestro-cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Maestro-cli version:", metadata.Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
