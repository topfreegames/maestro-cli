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
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login using an authorized email domain",
	Long: `Login and allow Maestro to authenticate commands using an authorized
	email domain. Google Oauth2 is used as authentication.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: inform server url")
			os.Exit(1)
		}
		log := newLog("login")

		url := args[0]
		log.Debugf("saving remote %s on config", url)

		home, err := homeDir()
		if err != nil {
			log.WithError(err).Fatal("getting home directory")
		}

		dirPath := filepath.Join(home, ".maestro")
		configPath := filepath.Join(dirPath, "config.yaml")

		config := make(map[string]interface{})
		if _, err := os.Stat(configPath); err == nil {
			log.Debug("reading file", configPath)
			bts, err := ioutil.ReadFile(configPath)
			if err != nil {
				log.WithError(err).Fatal("error reading file", configPath)
			}

			log.Debug("unmarshaling file", configPath)
			err = yaml.Unmarshal(bts, &config)
			if err != nil {
				log.WithError(err).Fatal("error unmarshaling file", configPath)
			}
		}

		config["url"] = url
		bts, err := yaml.Marshal(config)
		if err != nil {
			log.WithError(err).Fatalf("error marshaling obj: %#v", config)
		}

		log.Debug("mkdir ", dirPath)
		os.MkdirAll(dirPath, os.ModePerm)

		log.Debug("writing to file ", configPath)
		ioutil.WriteFile(configPath, bts, 0644)
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
