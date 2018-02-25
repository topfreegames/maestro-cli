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
	jsonLib "encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

type schedulerListResponse struct {
	Schedulers []string `json:"schedulers"`
}

var version string

// schedulerCmd represents the scheduler command
var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "list schedulers or get a specific one",
	Long:  `list schedulers or get a specific one`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("get scheduler")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		var url string
		if len(args) > 0 {
			schedulerName := args[0]
			url = fmt.Sprintf("%s/scheduler/%s/config?version=%s",
				config.ServerURL, schedulerName, version)
			body, status, err := client.Get(url)
			if err != nil {
				log.WithError(err).Fatal("error on get request")
			}

			if status != http.StatusOK {
				printError(body)
				return
			}

			obj := make(map[string]interface{})
			err = jsonLib.Unmarshal(body, &obj)
			if err != nil {
				log.WithError(err).Fatal("error on get request")
			}

			fmt.Println(obj["yaml"])
			return
		}

		url = fmt.Sprintf("%s/scheduler", config.ServerURL)
		body, status, err := client.Get(url)
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		if status != http.StatusOK {
			printError(body)
			return
		}

		var response *schedulerListResponse
		err = jsonLib.Unmarshal(body, &response)
		if err != nil {
			log.WithError(err).Fatal("error unmarshaling response")
		}

		for _, name := range response.Schedulers {
			fmt.Println(name)
		}
	},
}

func init() {
	getCmd.AddCommand(schedulerCmd)
	schedulerCmd.Flags().StringVar(&version, "version", "", "scheduler release version")
}
