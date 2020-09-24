// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2020 Wildlife Studios backend@tfgco.com

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Verbose determines how verbose maestro will run under
var Verbose int

var context string

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
		&Verbose, "verbose", "v", 0,
		"Verbosity level => v0: Error, v1=Warning, v2=Info, v3=Debug",
	)
	RootCmd.PersistentFlags().StringVarP(&context, "context", "c", "prod", "Maestro context, use it to manage different maestro clusters.")
}
