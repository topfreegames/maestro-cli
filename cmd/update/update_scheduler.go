// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package update

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

// updateSchedulerCmd represents the update command
var updateSchedulerCmd = &cobra.Command{
	Use:     "scheduler",
	Short:   "Updates scheduler",
	Example: "maestro-cli update scheduler ./scheduler.yaml",
	Long:    "Uses a .yaml file (argument) to update a new scheduler on Maestro.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewUpdateScheduler(client, config).run(cmd, args)
	},
}

type UpdateScheduler struct {
	client interfaces.Client
	config *extensions.Config
}

func NewUpdateScheduler(client interfaces.Client, config *extensions.Config) *UpdateScheduler {
	return &UpdateScheduler{
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

func (cs *UpdateScheduler) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()
	filePath := args[0]
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading scheduler file: %w", err)
	}
	if isYAML := common.IsYAML(filePath); !isYAML {
		return fmt.Errorf("file should be .yaml")
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

		operationId, err := cs.EnqueueUpdateSchedulerOperation(schedulerJsonBytes)
		if err != nil {
			return err
		}
		logger.Info("Successfully executed update scheduler. Operation id: " + operationId)
	}

	return nil
}

func (cs *UpdateScheduler) EnqueueUpdateSchedulerOperation(schedulerJsonBytes []byte) (string, error) {
	logger := common.GetLogger()

	var request v1.UpdateSchedulerRequest
	err := protojson.Unmarshal(schedulerJsonBytes, &request)
	if err != nil {
		return "", fmt.Errorf("error parsing Json to v1.UpdateSchedulerRequest: %w", err)
	}

	serializedRequest, err := common.Marshaller.Marshal(&request)
	if err != nil {
		return "", fmt.Errorf("error parsing request to json: %w", err)
	}

	logger.Debug("updating scheduler: " + request.Name)

	url := fmt.Sprintf("%s/schedulers", cs.config.ServerURL)

	body, status, err := cs.client.Put(url, string(serializedRequest))
	if err != nil {
		return "", fmt.Errorf("error on Put request: %w", err)
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("update scheduler response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	var response v1.UpdateSchedulerResponse
	err = protojson.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error deserializing update scheduler response, details: %w", err)
	}
	return response.OperationId, nil
}
