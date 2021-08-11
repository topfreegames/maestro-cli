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
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates new scheduler",
	Long: `Creates a new scheduler on Maestro and, if worker is running, the 
	rooms will be launghed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler config file path")
			os.Exit(1)
		}

		log := newLog("create")

		filePath := args[0]
		log.Debugf("reading %s", filePath)

		bts, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.WithError(err).Fatal("error while reading file")
		}

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		fmt.Println("Creating scheduler, this may take a few minutes...")

		url := fmt.Sprintf("%s/scheduler", config.ServerURL)
		body, status, err := client.Post(url, string(bts))
		if err != nil {
			log.WithError(err).Fatal("error on post request")
		}
		if status != http.StatusCreated {
			printError(body)
			return
		}

		fmt.Println("Successfully created scheduler")
		fmt.Println(string(bts))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
