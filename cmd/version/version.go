// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package version

import (
	"fmt"

	"github.com/topfreegames/maestro-cli/metadata"

	"github.com/spf13/cobra"
)

// Cmd represents the version command
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Maestro-cli version",
	Long:  `Maestro-cli version`,
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println("Maestro-cli version:", metadata.Version)
	},
}
