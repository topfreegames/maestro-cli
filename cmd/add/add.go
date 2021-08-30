// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package add

import (
	"github.com/spf13/cobra"
)

// Cmd represents the add command
var Cmd = &cobra.Command{
	Use:   "add",
	Short: "Addition operation",
	Long:  `Adds a resource, to know more type maestro-cli add --help.`,
}

func init() {
	Cmd.AddCommand(addRoomsCmd)
}
