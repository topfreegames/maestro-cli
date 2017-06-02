// maestro-cli api
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/maestro-cli/interfaces"
)

type Login struct {
	OAuthState string
	ServerURL  string
	OpenBroser func(string) error
}

func NewLogin(serverURL string, openBrower func(string) error) *Login {
	if openBrower == nil {
		openBrower = openBrowserFunc
	}
	return &Login{
		OAuthState: uuid.NewV4().String(),
		ServerURL:  serverURL,
		OpenBroser: openBrower,
	}
}

func (l *Login) Perform(client interfaces.Client) error {
	path := fmt.Sprintf("%s/login?state=%s", l.ServerURL, l.OAuthState)
	body, status, err := client.Get(path)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("status code %d when GET request to controller server", status)
	}
	var bodyObj map[string]string
	json.Unmarshal(body, &bodyObj)
	url := bodyObj["url"]

	err = l.OpenBroser(url)
	if err != nil {
		return err
	}
	return nil
}

func openBrowserFunc(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
