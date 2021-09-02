// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package remove

import (
	"github.com/spf13/cobra"
)

// Cmd represents the remove command
var Cmd = &cobra.Command{
	Use:   "remove",
	Short: "Removal operation",
	Long:  `Removes a resource, to know more, type maestro-cli remove --help.`,
}

func init() {
	Cmd.AddCommand(removeRoomsCmd)
}
