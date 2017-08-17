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

var replicas int

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "sets the number of replicas",
	Long:  `scales the scheduler up or down to match the number of replicas specified.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler name")
			os.Exit(1)
		}

		schedulerName := args[0]

		if replicas < 0 {
			fmt.Println("Error: inform replicas >= 0")
			os.Exit(1)
		}

		log := newLog("scale")

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		log.Debugf("reading %s", schedulerName)

		fmt.Printf("Scaling scheduler '%s' to have %d rooms, this may take a few minutes...\n", schedulerName, replicas)

		reqBody := map[string]interface{}{"replicas": replicas}
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

		fmt.Printf("Successfully scaled scheduler '%s' to %d replicas\n", schedulerName, replicas)
	},
}

func init() {
	RootCmd.AddCommand(scaleCmd)
	scaleCmd.Flags().IntVarP(&replicas, "replicas", "r", -1, "Total replicas of rooms to have")
}
