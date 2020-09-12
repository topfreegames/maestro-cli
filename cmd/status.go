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
	jsonLib "encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status SCHEDULER_NAME",
	Short: "Scheduler status",
	Long:  `Returns scheduler status like state, how many rooms are running and last time it scaled.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("inform scheduler name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("status")

		schedulerName := args[0]
		log.Debugf("reading %s", schedulerName)

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		url := fmt.Sprintf("%s/scheduler/%s", config.ServerURL, schedulerName)

		body, status, err := client.Get(url, "")
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		if status != http.StatusOK {
			printError(body)
			return
		}

		jsonBody := make(map[string]interface{})
		err = jsonLib.Unmarshal(body, &jsonBody)
		if err != nil {
			log.WithError(err).Fatal("error unmarshaling response")
		}

		bts, err := jsonLib.MarshalIndent(jsonBody, "", "\t")
		if err != nil {
			log.WithError(err).Fatal("error parsing response")
		}
		fmt.Println(string(bts))
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
