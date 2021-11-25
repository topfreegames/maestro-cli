// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package update

import (
	"github.com/spf13/cobra"
)

// Cmd represents the create command
var Cmd = &cobra.Command{
	Use:   "update",
	Short: "Updates a resource",
	Long:  `Updates a resource. To know more type maestro-cli update --help.`,
}

func init() {
	Cmd.AddCommand(updateSchedulerCmd)
}
