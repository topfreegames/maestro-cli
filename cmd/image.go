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
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var imageName string
var container string

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "sets a new image",
	Long: `updates scheduler with new image, changing only the scheduler's image field. If the image is the same,
	nothing is done. The rooms receive a gracefully shutdown and new ones are created.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform scheduler name")
			os.Exit(1)
		}

		log := newLog("setImage")

		schedulerName := args[0]
		if schedulerName == "" {
			fmt.Println("Error: inform scheduler name")
			os.Exit(1)
		}
		log.Debugf("updating %s", schedulerName)

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}
		client := getClient(config)

		fmt.Println("Updating scheduler image, this may take a few minutes...")
		url := fmt.Sprintf("%s/scheduler/%s/image?maxsurge=%s", config.ServerURL, schedulerName, maxsurge)
		reqBody := map[string]interface{}{"image": imageName, "container": container}
		reqBts, _ := json.Marshal(reqBody)
		body, status, err := client.Put(url, string(reqBts))
		if err != nil {
			log.WithError(err).Fatal("error on put request")
		}
		if status != http.StatusOK {
			printError(body)
			return
		}

		fmt.Printf("Successfully updated scheduler '%s' to image '%s'\n", schedulerName, imageName)
	},
}

func init() {
	setCmd.AddCommand(imageCmd)
	imageCmd.Flags().StringVarP(&imageName, "image", "i", "", "new image name")
	imageCmd.Flags().StringVar(&container, "container", "", "container name")

	imageCmd.Flags().StringVarP(&maxsurge, "maxsurge", "m", "", "percentage of the rooms to update at each step. Default is 25%.")
}
