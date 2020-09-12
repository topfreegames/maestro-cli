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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel OPERATION_KEY",
	Short: "Cancel an operation",
	Long:  `The operation will stop and rollback`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("specify an operation key")
		}

		operationKey := args[0]
		splitted := strings.Split(operationKey, ":")
		if len(splitted) < 2 {
			return errors.New("invalid operation key, it should be three values concatenated with a colon (:)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("cancel")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}

		client := getClient(config)

		// Get config from server
		operationKey := args[0]
		splitted := strings.Split(operationKey, ":")

		schedulerName := splitted[1]

		url := fmt.Sprintf("%s/scheduler/%s/operations/%s/cancel", config.ServerURL, schedulerName, operationKey)
		body, status, err := client.Put(url, "")
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		if status != http.StatusOK {
			log.Fatal(err)
		}

		fmt.Println(string(body))
	},
}

func init() {
	RootCmd.AddCommand(cancelCmd)
}
