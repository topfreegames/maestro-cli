// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"os"

	"github.com/topfreegames/maestro-cli/cmd/remove"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/cmd/add"
	"github.com/topfreegames/maestro-cli/cmd/create"
	"github.com/topfreegames/maestro-cli/cmd/get"
	initPkg "github.com/topfreegames/maestro-cli/cmd/init"
	"github.com/topfreegames/maestro-cli/cmd/version"
	"github.com/topfreegames/maestro-cli/common"
)

// RootCmd is the root command for maestro CLI application
var RootCmd = &cobra.Command{
	Use:   "maestro-cli",
	Short: "maestro-cli calls maestro api routes",
	Long:  `Use maestro-cli to control game rooms schedulers on Kubernetes.`,
}

// Execute runs RootCmd to initialize maestro CLI application
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(
		&common.Verbose, "verbose", "v", -1,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug",
	)
	RootCmd.PersistentFlags().StringVarP(&common.Context, "context", "c", "prod", "Maestro context, use it to manage different maestro clusters.")
	RootCmd.AddCommand(add.Cmd)
	RootCmd.AddCommand(remove.Cmd)
	RootCmd.AddCommand(initPkg.Cmd)
	RootCmd.AddCommand(create.Cmd)
	RootCmd.AddCommand(version.Cmd)
	RootCmd.AddCommand(get.Cmd)
}
