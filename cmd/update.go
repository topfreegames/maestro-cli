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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var maxsurge string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update SCHEDULER_CONFIG_FILE",
	Short: "Update scheduler on Maestro",
	Long: `Update scheduler on Maestro will update config on databases and, 
	if necessary, delete and create pods and services following new configuration.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("inform scheduler config file path")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("update")

		filePath := args[0]
		log.Debugf("reading %s", filePath)

		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.WithError(err).Fatal("error reading scheduler config")
		}

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		listBytes := bytes.Split(file, []byte("---"))
		for _, bts := range listBytes {
			if strings.TrimSpace(string(bts)) == "" {
				continue
			}

			yamlFile := make(map[string]interface{})
			err = yaml.Unmarshal(bts, &yamlFile)
			if err != nil {
				log.WithError(err).Fatal("error reading scheduler config")
			}
			schedulerName := yamlFile["name"].(string)
			url := fmt.Sprintf("%s/scheduler/%s?async=true&maxsurge=%s", config.ServerURL, schedulerName, maxsurge)
			body, status, err := client.Put(url, string(bts))
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

			fmt.Printf("Updating scheduler '%s', this may take a few minutes...\n", schedulerName)
			fmt.Printf("\nOperationKey\n===========\n%s\n", response["operationKey"])
			success := waitProgress(client, config, log, response["operationKey"].(string))
			if success {
				fmt.Printf("New scheduler config\n===================\n")
				fmt.Println(string(file))
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&maxsurge, "maxsurge", "m", "", "percentage of the rooms to update at each step. Default is 25%.")
}
