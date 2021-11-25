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
	"fmt"
	"io"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	yaml "gopkg.in/yaml.v2"
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
	var ll zapcore.Level
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

func IsYAML(path string) bool {
	return filepath.Ext(path) == ".yaml"
}

var Marshaller = &runtime.HTTPBodyMarshaler{
	Marshaler: &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	},
}
