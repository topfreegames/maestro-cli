package cmd

import (
	"errors"

	"github.com/topfreegames/maestro-cli/extensions"
)

func getConfig() (*extensions.Config, error) {
	filesystem := extensions.NewFileSystem()
	config, err := extensions.ReadConfig(filesystem)
	if err != nil {
		return nil, errors.New("probably you should login")
	}
	return config, nil
}

func getClient(config *extensions.Config) *extensions.Client {
	client := extensions.NewClient(config)
	return client
}
