// Copyright © 2018 TopFreeGames backend@tfgco.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback SCHEDULER_NAME VERSION",
	Short: "rollback to a previous version",
	Long:  `rollback to a previous version`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("rollback")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		var url string

		schedulerName := args[0]
		version := args[1]

		url = fmt.Sprintf("%s/scheduler/%s/rollback", config.ServerURL, schedulerName)
		reqBody := fmt.Sprintf(`{"version": "%s"}`, version)
		body, status, err := client.Put(url, reqBody)
		if err != nil {
			log.WithError(err).Fatal("error on put request")
		}

		if status != http.StatusOK {
			printError(body)
			return
		}

		fmt.Println(string(body))
	},
}

func init() {
	RootCmd.AddCommand(rollbackCmd)
}
