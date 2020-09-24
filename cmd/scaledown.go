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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// scaledownCmd represents the scaledown command
var scaledownCmd = &cobra.Command{
	Use:   "scaledown SCHEDULER_NAME",
	Short: "Scales down rooms at a specific scheduler",
	Long:  `Scales down rooms at a specific scheduler`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("inform scheduler name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		schedulerName := args[0]

		if amount == 0 {
			fmt.Println("Error: inform amount > 0")
			os.Exit(1)
		}

		log := newLog("scaledown")

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		log.Debugf("reading %s", schedulerName)

		fmt.Printf("Scaling down scheduler '%s' in %d rooms, this may take a few minutes...\n", schedulerName, amount)

		reqBody := map[string]interface{}{"scaledown": amount}
		reqBts, _ := json.Marshal(reqBody)
		url := fmt.Sprintf("%s/scheduler/%s", config.ServerURL, schedulerName)
		body, status, err := client.Post(url, string(reqBts))
		if err != nil {
			log.WithError(err).Fatal("error on post request")
		}
		if status != http.StatusOK {
			printError(body)
			return
		}

		fmt.Println("Successfully scaled down scheduler", schedulerName)
	},
}

func init() {
	RootCmd.AddCommand(scaledownCmd)
	scaledownCmd.Flags().UintVarP(&amount, "amount", "a", 0, "Amount of rooms to scale down")
}
