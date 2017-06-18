// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cmd

import (
	"errors"

	"github.com/topfreegames/maestro-cli/extensions"
)

func getConfig() (*extensions.Config, error) {
	filesystem := extensions.NewFileSystem()
	config, err := extensions.ReadConfig(filesystem, context)
	if err != nil {
		return nil, errors.New("probably you should login")
	}
	return config, nil
}

func getClient(config *extensions.Config) *extensions.Client {
	client := extensions.NewClient(config)
	return client
}
