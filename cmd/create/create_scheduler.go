// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package create

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	"google.golang.org/protobuf/encoding/protojson"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	yaml "gopkg.in/yaml.v2"
)

var marshler = &runtime.HTTPBodyMarshaler{
	Marshaler: &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	},
}

// createSchedulerCmd represents the create command
var createSchedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Creates new scheduler",
	Long:  `Uses a file (argument) to create a new scheduler on Maestro`,
	Args:  validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := getClientAndConfig()
		if err != nil {
			return err
		}

		return NewCreateScheduler(client, config).run(cmd, args)
	},
}

type CreateScheduler struct {
	client interfaces.Client
	config *extensions.Config
}

func NewCreateScheduler(client interfaces.Client, config *extensions.Config) *CreateScheduler {
	return &CreateScheduler{
		client: client,
		config: config,
	}
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("error: missing arg with scheduler config file path")
	}

	filePath := args[0]
	_, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	return nil
}

func (cs *CreateScheduler) run(_ *cobra.Command, args []string) error {

	filePath := args[0]
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading scheduler file: %w", err)
	}

	r := bytes.NewReader(bts)
	dec := yaml.NewDecoder(r)

	for {
		var request v1.CreateSchedulerRequest
		decodeErr := dec.Decode(&request)
		if decodeErr != nil {
			if decodeErr.Error() == "EOF" {
				return nil
			}
			return decodeErr
		}

		serializedRequest, err := marshler.Marshal(&request)
		if err != nil {
			return fmt.Errorf("error parsing request to json: %w", err)
		}

		fmt.Println("creating scheduler: ", request.Name)

		url := fmt.Sprintf("%s/schedulers", cs.config.ServerURL)
		body, status, err := cs.client.Post(url, string(serializedRequest))
		if err != nil {
			return fmt.Errorf("error on post request: %w", err)
		}
		if status != http.StatusOK {
			return fmt.Errorf("create scheduler response not ok, status: %s, body: %s", http.StatusText(status), string(body))
		}

		fmt.Println("Successfully created scheduler: ", request.Name)
	}
}

func getClientAndConfig() (interfaces.Client, *extensions.Config, error) {

	config, err := common.GetConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error getting client config: %w", err)
	}

	client := common.GetClient(config)

	return client, config, nil
}
