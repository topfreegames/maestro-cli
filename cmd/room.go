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
	"github.com/spf13/cobra"
	"net/http"
)

// roomsCmd represents the rooms command
var roomCmd = &cobra.Command{
	Use:   "room SCHEDULER_NAME ROOM_ID",
	Short: "show the game rooms details",
	Long:  `show the game room details of a scheduler`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("room")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		var url string

		if len(args) < 2 {
			log.Fatal("error: you must specify scheduler name and room id")
			return
		}
		schedulerName := args[0]
		roomId := args[1]

		url = fmt.Sprintf("%s/scheduler/%s/rooms/%s", config.ServerURL, schedulerName, roomId)
		body, responseStatus, err := client.Get(url, "")
		if err != nil {
			log.WithError(err).Fatal("error on get room request")
		}

		if responseStatus != http.StatusOK {
			printError(body)
			return
		}

		var room gameRoom

		err = json.Unmarshal(body, &room)
		if err != nil {
			log.WithError(err).Fatal("error deserializing get rooms response body")
		}

		printRoomsTable([]gameRoom{room})
	},
}

func init() {
	RootCmd.AddCommand(roomCmd)
}
