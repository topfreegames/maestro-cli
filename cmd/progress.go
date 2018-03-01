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
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/extensions"
)

// progressCmd represents the progress command
var progressCmd = &cobra.Command{
	Use:   "progress OPERATION_KEY",
	Short: "Returns if the operation of operationKey is enqueued or the progress of the operation",
	Long:  `Returns if the operation of operationKey is enqueued or the progress of the operation`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("progress")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		if len(args) == 0 {
			log.Fatal("error: specify scheduler name")
			return
		}

		// Get config from server
		operationKey := args[0]

		waitProgress(client, config, log, operationKey)
	},
}

func waitProgress(client *extensions.Client, config *extensions.Config, log *logrus.Logger, operationKey string) bool {
	splitted := strings.Split(operationKey, ":")
	if len(splitted) < 2 {
		log.Fatal("error: invalid operation key")
		return false
	}

	schedulerName := splitted[1]

	bars := []string{"|", "/", "-", "\\"}
	i := 0

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			url := fmt.Sprintf("%s/scheduler/%s/operations/%s/status", config.ServerURL, schedulerName, operationKey)
			body, status, err := client.Get(url, "")
			if err != nil {
				log.WithError(err).Fatal("error on get request")
			}

			if status != http.StatusOK {
				fmt.Println("")
				fmt.Println(string(body))
				return false
			}

			var response map[string]interface{}
			json.Unmarshal(body, &response)

			if _, ok := response["success"]; ok {
				fmt.Println("")
				fmt.Println(string(body))
				return true
			}

			fmt.Printf("\r[%s] %s %s", bars[i], response["operation"], response["progress"])
			i = (i + 1) % len(bars)
		}
	}

	return true
}

func init() {
	RootCmd.AddCommand(progressCmd)
}
