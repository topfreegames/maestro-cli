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
	"net/url"

	"github.com/topfreegames/maestro-cli/api"
	"github.com/topfreegames/maestro-cli/extensions"
	loginExt "github.com/topfreegames/maestro-cli/login"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [maestro api url]",
	Short: "Login using an authorized email domain",
	Long: `Login and allow Maestro to authenticate commands using an authorized
	email domain. Google Oauth2 is used as authentication.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("inform server url")
		}

		maestroURL := args[0]
		_, err := url.ParseRequestURI(maestroURL)

		if err != nil {
			return errors.New("inform server url in a valid format")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log := newLog("login")

		url := args[0]
		log.Debugf("saving remote %s on config", url)

		login := loginExt.NewLogin(url, nil)
		client := extensions.NewClient(nil)
		filesystem := extensions.NewFileSystem()

		app, err := api.NewApp(login, filesystem, client, log, context)
		if err != nil {
			log.WithError(err).Fatal("error with app constructor")
		}

		closer, err := app.ListenAndLoginAndServe()
		if err != nil {
			log.WithError(err).Fatal("error with app constructor")
		}
		closer.Close()
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
