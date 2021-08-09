// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package common

import (
	"bytes"
	"errors"

	"github.com/topfreegames/maestro-cli/extensions"
	yaml "gopkg.in/yaml.v2"
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

func SplitYAML(resources []byte) ([][]byte, error) {

	dec := yaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := yaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}
