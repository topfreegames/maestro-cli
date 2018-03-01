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
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("cancel")
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
		splitted := strings.Split(operationKey, ":")
		if len(splitted) < 2 {
			log.Fatal("error: invalid operation key")
			return
		}

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cancelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cancelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
