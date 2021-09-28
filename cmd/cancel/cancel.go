// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cancel

import (
	"github.com/spf13/cobra"
)

// Cmd represents the cancel command
var Cmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel operation",
	Long:  `Cancels a resource, to know more, type maestro-cli cancel --help.`,
}

func init() {
	Cmd.AddCommand(cancelOperationCmd)
}
