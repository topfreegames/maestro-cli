// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create scheduler",
	Short: "Creates new scheduler",
	Long: `Creates a new scheduler on Maestro and, if worker is running, the 
	rooms will be launghed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler config file path")
			os.Exit(1)
		}

		filePath := args[0]
		zap.L().Sugar().Debugf("reading %s", filePath)

		bts, err := ioutil.ReadFile(filePath)
		if err != nil {
			zap.L().With(zap.Error(err)).Fatal("error while reading file")
			os.Exit(1)
		}

		config, err := getConfig()
		if err != nil {
			zap.L().With(zap.Error(err)).Fatal("error getting client config")
			os.Exit(1)
		}
		client := getClient(config)

		fmt.Println("Creating scheduler...")

		url := fmt.Sprintf("%s/scheduler", config.ServerURL)
		body, status, err := client.Post(url, string(bts))
		if err != nil {
			zap.L().With(zap.Error(err)).Fatal("error on post request")
			os.Exit(1)
		}
		if status != http.StatusCreated {
			zap.L().Error("create scheduler response not ok", 
				zap.String("status", http.StatusText(status)),
				zap.String("body", string(body)))
			os.Exit(1)
		}

		fmt.Println("Successfully created scheduler")
		fmt.Println(string(bts))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
