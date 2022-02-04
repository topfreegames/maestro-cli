// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package create

import (
	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/cmd/create/scheduler"
	"github.com/topfreegames/maestro-cli/cmd/create/scheduler_version"
)

// Cmd represents the create command
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a resource",
	Long:  `Creates a resource, to know more type maestro-cli create --help.`,
}

func init() {
	Cmd.AddCommand(scheduler.CreateSchedulerCmd)
	Cmd.AddCommand(scheduler_version.CreateSchedulerVersionCmd)
}
