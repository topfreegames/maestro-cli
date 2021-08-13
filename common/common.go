// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package common

import (
	"errors"
	"fmt"

	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Verbose determines how verbose maestro will run under
var Verbose int

var Context string

func GetClientAndConfig() (interfaces.Client, *extensions.Config, error) {

	config, err := GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting client config: %w", err)
	}

	client := GetClient(config)

	return client, config, nil
}

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

func GetLogger() *zap.Logger {
	ll := zap.InfoLevel
	switch Verbose {
	case 0:
		ll = zap.ErrorLevel
	case 1:
		ll = zap.WarnLevel
	case 3:
		ll = zap.DebugLevel
	default:
		ll = zap.InfoLevel
	}

	log := zap.NewDevelopmentConfig()
	log.OutputPaths = []string{"stdout"}
	log.Level.SetLevel(ll)
	log.EncoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",
	}

	logger, _ := log.Build()

	return logger
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
