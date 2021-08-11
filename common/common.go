// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package common

import (
	"errors"

	"github.com/topfreegames/maestro-cli/extensions"
)

// Verbose determines how verbose maestro will run under
var Verbose int

var Context string

func GetConfig() (*extensions.Config, error) {
	filesystem := extensions.NewFileSystem()
	config, err := extensions.ReadConfig(filesystem, Context)
	if err != nil {
		return nil, errors.New("probably you should login")
	}
	return config, nil
}

func GetClient(config *extensions.Config) *extensions.Client {
	client := extensions.NewClient(config)
	return client
}

func Success(response map[string]interface{}) (string, bool) {
	if response["success"] == false {
		if reason, ok := response["reason"].(string); ok {
			return reason, false
		}

		return "failed", false
	}

	return "", true
}
