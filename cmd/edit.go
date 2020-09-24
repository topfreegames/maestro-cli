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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var editor string

func getEditor() string {
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	if editor == "" {
		editor = "vi"
	}

	return editor
}

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit SCHEDULER_NAME",
	Short: "Edit a scheduler",
	Long: `edit opens default editor and updates the scheduler on save if scheduler is valid.
	It chooses which editor to use by reading $EDITOR variable; if it's unset, maestro chooses vi as default.
	To specify which editor to use, use the --editor flag.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("inform scheduler name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("get scheduler")

		config, err := getConfig()
		if err != nil {
			log.WithError(err).Fatal("error getting client config")
		}

		client := getClient(config)

		// Get config from server
		schedulerName := args[0]
		url := fmt.Sprintf("%s/scheduler/%s/config", config.ServerURL, schedulerName)
		body, status, err := client.Get(url, "")
		if err != nil {
			log.WithError(err).Fatal("error on get request")
		}

		if status != http.StatusOK {
			log.Fatal(err)
		}

		var bodyJSON map[string]string
		err = json.Unmarshal(body, &bodyJSON)
		if err != nil {
			log.WithError(err).Fatal("error unmarshalling response from server")
		}

		yamlString := bodyJSON["yaml"]

		// Create tmp file
		fileName := fmt.Sprintf("/tmp/%s-%s.yaml", schedulerName, uuid.New().String())
		err = ioutil.WriteFile(fileName, []byte(yamlString), 0644)
		if err != nil {
			log.WithError(err).Fatal("error writing file on /tmp")
		}

		// Open on editor
		editor = getEditor()
		editorCmd := exec.Command(editor, fileName)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		err = editorCmd.Start()
		if err != nil {
			log.WithError(err).Fatal("error on start $EDITOR")
		}

		err = editorCmd.Wait()
		if err != nil {
			log.WithError(err).Fatal("error while editing")
		}

		// Read new file
		updatedYamlBts, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.WithError(err).Fatal("error reading saved file")
		}

		// Remove new file
		err = os.Remove(fileName)
		if err != nil {
			log.WithError(err).Fatal("error removing file from /tmp")
		}

		updatedYamlString := string(updatedYamlBts)
		if updatedYamlString == yamlString {
			fmt.Println("nothing changed")
			return
		}

		fmt.Printf("Updating scheduler %s. This can take a while...\n", schedulerName)

		// Update scheduler
		url = fmt.Sprintf("%s/scheduler/%s?async=true&maxsurge=%s", config.ServerURL, schedulerName, maxsurge)
		body, status, err = client.Put(url, updatedYamlString)
		if err != nil {
			log.WithError(err).Fatal("error on put request")
		}
		if status != http.StatusOK {
			printError(body)
			return
		}

		var response map[string]interface{}
		json.Unmarshal(body, &response)
		if reason, ok := success(response); !ok {
			fmt.Printf("Operation failed. Try again later.\nReason: %s\n", reason)
			return
		}

		fmt.Printf("\nOperationKey\n===========\n%s\n", response["operationKey"])
		success := waitProgress(client, config, log, response["operationKey"].(string))
		if success {
			fmt.Printf("New scheduler config\n===================\n")
			fmt.Println(updatedYamlString)
		}
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&maxsurge, "maxsurge", "m", "", "percentage of the rooms to update at each step. Default is 25%.")
	editCmd.Flags().StringVarP(&editor, "editor", "e", "", "specify which editor to use.")
}
