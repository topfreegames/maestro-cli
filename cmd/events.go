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
	"strings"
	"text/tabwriter"
	"time"
)

var page int

type schedulerEvent struct {
	Name          string
	SchedulerName string
	CreatedAt     string
	Metadata      map[string]interface{}
}

// eventsCmd represents the events command
var eventsCmd = &cobra.Command{
	Use:   "events SCHEDULER_NAME --page PAGE",
	Short: "list the events of a scheduler",
	Long:  `list the events of a scheduler, it uses pagination with page number default equal 1 and page size equal 30`,
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("events")
		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)
		var url string

		if len(args) == 0 {
			log.Fatal("error: you must specify scheduler name")
			return
		}
		schedulerName := args[0]

		url = fmt.Sprintf("%s/scheduler/%s/events?page=%d", config.ServerURL, schedulerName, page)
		body, status, err := client.Get(url, "")
		if err != nil {
			log.WithError(err).Fatal("error on get events request")
		}

		if status != http.StatusOK {
			printError(body)
			return
		}

		var events []schedulerEvent

		err = json.Unmarshal(body, &events)
		if err != nil {
			log.WithError(err).Fatal("error deserializing get events response body")
		}

		printEventsTable(events)
	},
}

func printEventsTable(events []schedulerEvent) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	format := "%s\t\t%s\t\t%s\t\n"
	fmt.Fprintf(w, format, "EVENT_NAME", "AGE", "METADATA")

	for _, event := range events {
		createdAt, _ := time.Parse(time.RFC3339Nano, event.CreatedAt)
		age := time.Since(createdAt)
		prettyAge := durafmt.ParseShort(age).String()
		prettyMetadata, _ := json.Marshal(event.Metadata)

		fmt.Fprintf(
			w,
			format,
			strings.ToUpper(event.Name),
			prettyAge,
			string(prettyMetadata),
		)
	}
}

func init() {
	RootCmd.AddCommand(eventsCmd)
	eventsCmd.Flags().IntVarP(&page, "page", "p", 1, "Scheduler events pagination number")
}
