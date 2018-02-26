// maestro-cli api
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
)

func SaveAccessToken(
	state, code, expectedState, serverURL,
	context string,
	fs interfaces.FileSystem,
	client interfaces.Client,
) error {
	if state != expectedState {
		err := fmt.Errorf("invalid oauth state, expected '%s', got '%s'", expectedState, state)
		return err
	}

	url := fmt.Sprintf("%s/access?code=%s", serverURL, code)
	resp, status, err := client.Get(url, "")
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return errors.New(string(resp))
	}

	var bodyObj map[string]interface{}
	json.Unmarshal(resp, &bodyObj)
	token := bodyObj["token"].(string)

	c := extensions.NewConfig(token, serverURL)
	err = c.Write(fs, context)
	return err
}
