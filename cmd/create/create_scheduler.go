// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package create

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	"google.golang.org/protobuf/encoding/protojson"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	k8s_yaml "sigs.k8s.io/yaml"
)

// createSchedulerCmd represents the create command
var createSchedulerCmd = &cobra.Command{
	Use:     "scheduler",
	Short:   "Creates new scheduler",
	Example: "maestro-cli create scheduler ./scheduler.yaml",
	Long:    "Uses a file (argument) to create a new scheduler on Maestro.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
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
		return errors.New("missing arg with scheduler config file path")
	}

	filePath := args[0]
	_, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	return nil
}

func (cs *CreateScheduler) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()
	filePath := args[0]
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading scheduler file: %w", err)
	}

	yamls, err := common.SplitYAML(bts)
	if err != nil {
		return fmt.Errorf("error splitting YAML file into multiple objects: %w", err)
	}

	for _, yaml_object := range yamls {

		schedulerJsonBytes, err := k8s_yaml.YAMLToJSON(yaml_object)
		if err != nil {
			return fmt.Errorf("error parsing YAML to Json: %w", err)
		}

		var request v1.CreateSchedulerRequest
		err = protojson.Unmarshal(schedulerJsonBytes, &request)
		if err != nil {
			return fmt.Errorf("error parsing Json to v1.CreateSchedulerRequest: %w", err)
		}

		serializedRequest, err := common.Marshaller.Marshal(&request)
		if err != nil {
			return fmt.Errorf("error parsing request to json: %w", err)
		}

		logger.Debug("creating scheduler: " + request.Name)

		url := fmt.Sprintf("%s/schedulers", cs.config.ServerURL)
		body, status, err := cs.client.Post(url, string(serializedRequest))
		if err != nil {
			return fmt.Errorf("error on post request: %w", err)
		}
		if status != http.StatusOK {
			return fmt.Errorf("create scheduler response not ok, status: %s, body: %s", http.StatusText(status), string(body))
		}

		logger.Info("Successfully created scheduler: " + request.Name)
	}

	return nil
}
