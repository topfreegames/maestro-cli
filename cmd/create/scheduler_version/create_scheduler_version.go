package scheduler_version

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"

	k8s_yaml "sigs.k8s.io/yaml"
)

var CreateSchedulerVersionCmd = &cobra.Command{
	Use:     "scheduler-version",
	Short:   "Creates new scheduler version",
	Example: "maestro-cli create scheduler-version ./scheduler.yaml",
	Long:    "Uses a .yaml file (argument) to create new scheduler version on Maestro.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewCreateSchedulerVersion(client, config).run(cmd, args)
	},
}

type CreateSchedulerVersion struct {
	client interfaces.Client
	config *extensions.Config
}

func NewCreateSchedulerVersion(client interfaces.Client, config *extensions.Config) *CreateSchedulerVersion {
	return &CreateSchedulerVersion{
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

func (cs *CreateSchedulerVersion) run(_ *cobra.Command, args []string) error {
	logger := common.GetLogger()
	filePath := args[0]
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading scheduler version file: %w", err)
	}
	if isYaml := common.IsYAML(filePath); !isYaml {
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

		operationId, err := cs.EnqueueNewSchedulerVersionOperation(schedulerJsonBytes)
		if err != nil {
			return err
		}
		logger.Info("Successfully executed new scheduler version. Operation id: " + operationId)
	}

	return nil
}

func (cs *CreateSchedulerVersion) EnqueueNewSchedulerVersionOperation(schedulerJsonBytes []byte) (string, error) {
	logger := common.GetLogger()
	var request v1.NewSchedulerVersionRequest
	err := protojson.Unmarshal(schedulerJsonBytes, &request)
	if err != nil {
		return "", fmt.Errorf("error parsing Json to v1.NewSchedulerVersionRequest: %w", err)
	}

	serializedRequest, err := common.Marshaller.Marshal(&request)
	if err != nil {
		return "", fmt.Errorf("error parsing request to json: %w", err)
	}

	logger.Debug("updating scheduler: " + request.Name)

	url := fmt.Sprintf("%s/schedulers/%s", cs.config.ServerURL, request.Name)

	body, status, err := cs.client.Post(url, string(serializedRequest))
	if err != nil {
		return "", fmt.Errorf("error on Post request: %w", err)
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("new scheduler version response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	var response v1.NewSchedulerVersionResponse
	err = protojson.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error deserializing new scheduler version response, details: %w", err)
	}
	return response.OperationId, nil
}
