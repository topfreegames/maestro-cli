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
	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

const (
	defaultLimit = 30

	//StatusCreating represents a room status
	StatusCreating = "creating"

	//StatusReady represents a room status
	StatusReady = "ready"

	//StatusOccupied represents a room status
	StatusOccupied = "occupied"

	//StatusReadyOrOccupied represents an aggregate of room status
	StatusReadyOrOccupied = "ready_or_occupied"

	//StatusTerminating represents a room status
	StatusTerminating = "terminating"

	//StatusTerminated represents a room status
	StatusTerminated = "terminated"
)

type gameRoom struct {
	RoomId           string
	SchedulerName    string
	SchedulerVersion string
	Status           string
	CreatedAt        string
	LastPingAt       string
}

// roomsCmd represents the rooms command
var roomsCmd = &cobra.Command{
	Use:   "rooms SCHEDULER_NAME STATUS --page PAGE",
	Short: "list the game rooms of a scheduler",
	Long:  `list the game rooms of a scheduler, it uses pagination with page number default equal 1 and page size equal 30`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("rooms")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		var url string

		if len(args) < 2 {
			log.Fatal("error: you must specify scheduler name and status")
			return
		}
		schedulerName := args[0]
		status := args[1]

		if !isValidStatus(status) {
			log.Fatal(fmt.Sprintf(
				"error: %s is an invalid status option, valid options are: [%s, %s, %s, %s, %s, %s]",
				status,
				StatusTerminated,
				StatusTerminating,
				StatusReady,
				StatusOccupied,
				StatusReadyOrOccupied,
				StatusCreating,
			))
			return
		}

		url = fmt.Sprintf("%s/scheduler/%s/rooms/status/%s?limit=%d&offset=%d", config.ServerURL, schedulerName, status, defaultLimit, page)
		body, responseStatus, err := client.Get(url, "")
		if err != nil {
			log.WithError(err).Fatal("error on get room request")
		}

		if responseStatus != http.StatusOK {
			printError(body)
			return
		}

		var rooms []gameRoom

		err = json.Unmarshal(body, &rooms)
		if err != nil {
			log.WithError(err).Fatal("error deserializing get rooms response body")
		}

		printRoomsTable(rooms)
	},
}

func printRoomsTable(rooms []gameRoom) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	format := "%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\n"
	fmt.Fprintf(w, format, "SCHEDULER_NAME", "SCHEDULER_VERSION", "ROOM_ID", "STATUS", "ROOM_AGE", "LAST_PING_AGE")

	for _, room := range rooms {

		createdAt, _ := time.Parse(time.RFC3339Nano, room.CreatedAt)
		prettyAge := durafmt.ParseShort(time.Since(createdAt)).String()
		lastPingAt, _ := time.Parse(time.RFC3339Nano, room.LastPingAt)
		pingPrettyAge := durafmt.ParseShort(time.Since(lastPingAt)).String()
		if lastPingAt.Before(createdAt) || lastPingAt.Equal(createdAt) {
			pingPrettyAge = "no ping sent yet"
		}

		fmt.Fprintf(
			w,
			format,
			room.SchedulerName,
			room.SchedulerVersion,
			room.RoomId,
			room.Status,
			prettyAge,
			pingPrettyAge,
		)
	}
}

func isValidStatus(status string) bool {
	switch status {
	case StatusCreating, StatusReady, StatusOccupied, StatusTerminating, StatusTerminated, StatusReadyOrOccupied:
		return true
	}
	return false
}

func init() {
	RootCmd.AddCommand(roomsCmd)
	roomsCmd.Flags().IntVarP(&page, "page", "p", 1, "Scheduler rooms pagination number")
}
