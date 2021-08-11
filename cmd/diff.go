// Copyright © 2018 TFGco backend@tfgco.com
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

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff SCHEDULER_NAME [VERSION_1] [VERSION_2]",
	Short: "diff between configs of a scheduler",
	Long: `returns the diff between two versions of a scheduler. 
If no VERSION_1 is specified, VERSION_1 defaults to current version and VERSION_2 to the one before that.
If only VERSION_1 is specified, VERSION_2 defaults to the one before that.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("releases")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		if len(args) == 0 {
			fmt.Println("error: specify scheduler name")
			return
		}

		var version1, version2 string
		if len(args) > 1 {
			version1 = args[1]
		}

		if len(args) > 2 {
			version2 = args[2]
		}

		schedulerName := args[0]
		url := fmt.Sprintf("%s/scheduler/%s/diff", config.ServerURL, schedulerName)
		reqBody := fmt.Sprintf(`{"version1": "%s", "version2": "%s"}`, version1, version2)
		body, status, err := client.Get(url, reqBody)
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		if status != http.StatusOK {
			printError(body)
			return
		}

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.WithError(err).Fatal("error unmarshaling response")
		}

		title := fmt.Sprintf("%s: diff between %s and %s", schedulerName, response["version1"], response["version2"])
		bar := buildBar(title)
		fmt.Printf("%s\n%s\n", title, bar)
		fmt.Println(response["diff"])
	},
}

func init() {
	RootCmd.AddCommand(diffCmd)
}
