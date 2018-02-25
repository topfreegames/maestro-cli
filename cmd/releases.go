// Copyright © 2018 TFGCo backend@tfgco.com
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

	"github.com/spf13/cobra"
)

// releasesCmd represents the releases command
var releasesCmd = &cobra.Command{
	Use:   "releases",
	Short: "list the releases of a scheduler",
	Long:  `list the releases of a scheduler`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("releases")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		var url string

		schedulerName := args[0]
		url = fmt.Sprintf("%s/scheduler/%s/releases", config.ServerURL, schedulerName)
		body, status, err := client.Get(url)
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		if status != http.StatusOK {
			printError(body)
			return
		}

		obj := make(map[string]interface{})
		err = json.Unmarshal(body, &obj)
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		releases := obj["releases"].([]interface{})

		title := fmt.Sprintf("%s releases", schedulerName)
		l := len(title)
		bar := ""
		for i := 0; i < l; i++ {
			bar = bar + "="
		}

		fmt.Printf("%s\n%s\n", title, bar)
		for _, release := range releases {
			mapRelease := release.(map[string]interface{})
			fmt.Printf("%s\t%s\n", mapRelease["version"], mapRelease["createdAt"])
		}

		return
	},
}

func init() {
	RootCmd.AddCommand(releasesCmd)
}
