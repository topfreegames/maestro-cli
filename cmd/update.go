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
	"fmt"
	"io/ioutil"
	"os"

	jsonLib "encoding/json"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update scheduler on Maestro",
	Long: `Update scheduler on Maestro will update config on databases and, 
	if necessary, delete and create pods and services following new configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler config file path")
			os.Exit(1)
		}

		log := newLog("update")

		filePath := args[0]
		log.Debugf("reading %s", filePath)

		bts, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.WithError(err).Fatal("error reading scheduler config")
		}

		scheduler := make(map[string]interface{})
		err = jsonLib.Unmarshal(bts, &scheduler)
		if err != nil {
			log.WithError(err).Fatal("error unmarshaling scheduler config")
		}
		schedulerName, ok := scheduler["name"].(string)
		if !ok {
			log.WithError(err).Fatal("scheduler name should be a string")
		}

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		url := fmt.Sprintf("%s/scheduler/%s", config.ServerURL, schedulerName)

		body, status, err := client.Put(url, scheduler)
		if err != nil {
			log.WithError(err).Fatal("error on put request")
		}

		fmt.Println("Status:", status)
		fmt.Println("Response:", string(body))
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
