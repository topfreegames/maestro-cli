// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package get

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a resource",
	Long:  `Get a resource, to know more type maestro-cli list --help.`,
}

func init() {
	Cmd.AddCommand(getSchedulersCmd)
	Cmd.AddCommand(getOperationsCmd)
	Cmd.AddCommand(getSchedulersInfoCmd)
	Cmd.AddCommand(getOperationCmd)
}
