// Copyright © 2017 TopFreeGames backend@tfgco.com
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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var schedulerMin uint

// minCmd represents the min command
var minCmd = &cobra.Command{
	Use:   "min",
	Short: "sets scheduler's min",
	Long: `updates scheduler with new min, changing only the scheduler's min field. If the min is the same,
	nothing is done.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler name")
			os.Exit(1)
		}

		log := newLog("setImage")

		schedulerName := args[0]
		if schedulerName == "" {
			fmt.Println("Error: inform scheduler name")
			os.Exit(1)
		}

		log.Debugf("updating %s", schedulerName)

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		url := fmt.Sprintf("%s/scheduler/%s/min?async=true", config.ServerURL, schedulerName)
		reqBody := map[string]interface{}{"min": schedulerMin}
		reqBts, _ := json.Marshal(reqBody)
		body, status, err := client.Put(url, string(reqBts))
		if err != nil {
			log.WithError(err).Fatal("error on put request")
		}
		if status != http.StatusOK {
			printError(body)
			return
		}

		var response map[string]interface{}
		json.Unmarshal(body, &response)
		if reason, ok := success(response); !ok {
			fmt.Printf("Operation failed. Try again later.\nReason: %s\n", reason)
			return
		}

		fmt.Printf("Updating scheduler '%s' to min '%d'. This can take a few minutes...\n", schedulerName, schedulerMin)
		fmt.Printf("\nOperationKey\n===========\n%s\n", response["operationKey"])

		waitProgress(client, config, log, response["operationKey"].(string))
	},
}

func init() {
	setCmd.AddCommand(minCmd)
	minCmd.Flags().UintVarP(&schedulerMin, "min", "m", uint(0), "new scheduler min")
}
