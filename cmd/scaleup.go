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

var amount uint

// scaleupCmd represents the scaleup command
var scaleupCmd = &cobra.Command{
	Use:   "scaleup",
	Short: "Scales up rooms at a specific scheduler",
	Long:  `Scales up rooms at a specific scheduler`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler name")
			os.Exit(1)
		}

		schedulerName := args[0]

		if amount == 0 {
			fmt.Println("Error: inform amount > 0")
			os.Exit(1)
		}

		log := newLog("scaleup")

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		log.Debugf("reading %s", schedulerName)

		fmt.Printf("Scaling up scheduler '%s' in %d rooms, this may take a few minutes...\n", schedulerName, amount)

		reqBody := map[string]interface{}{"scaleup": amount}
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

		fmt.Println("Successfully scaled up scheduler", schedulerName)
	},
}

func init() {
	RootCmd.AddCommand(scaleupCmd)
	scaleupCmd.Flags().UintVarP(&amount, "amount", "a", 0, "Amount of rooms to scale up")
}
