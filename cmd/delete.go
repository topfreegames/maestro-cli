// Copyright © 2020 Wildlife Studios backend@tfgco.com
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
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [SCHEDULER_NAME...]",
	Short: "Delete scheduler on Maestro",
	Long: `Delete scheduler on Maestro. This will stop all rooms and remove
	configs from databases.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("inform scheduler name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("delete")

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		for _, schedulerName := range args {
			log.Debugf("reading %s", schedulerName)

			fmt.Printf("Deleting scheduler %s, this may take a few minutes...\n", schedulerName)

			url := fmt.Sprintf("%s/scheduler/%s", config.ServerURL, schedulerName)
			body, status, err := client.Delete(url)
			if err != nil {
				log.WithError(err).Fatal("error on delete request")
			}
			if status != http.StatusOK {
				printError(body)
				return
			}

			fmt.Println("Successfully deleted scheduler", schedulerName)
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
