package _switch

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
)

// switchActiveVersionCmd represents the switch command
var switchActiveVersionCmd = &cobra.Command{
	Use:     "active-version",
	Short:   "Switch active version to a given scheduler",
	Example: "maestro-cli switch active-version <scheduler_name> <target_version>",
	Long:    "Given the scheduler name and the target version, switch active version from Scheduler in Maestro-Next.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewSwitchActiveVersion(client, config).run(cmd, args)
	},
}

type SwitchActiveVersion struct {
	client interfaces.Client
	config *extensions.Config
}

func NewSwitchActiveVersion(client interfaces.Client, config *extensions.Config) *SwitchActiveVersion {
	return &SwitchActiveVersion{
		client: client,
		config: config,
	}
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("missing args: scheduler name or/and target version")
	}

	targetVersion := args[1]

	err := isVersionValid(targetVersion)
	if err != nil {
		return err
	}

	return nil
}

func isVersionValid(version string) error {
	_, err := semver.NewVersion(version)
	if err != nil {
		return errors.New("invalid version value. version must be in a semantic format")
	}
	return nil
}

func (a *SwitchActiveVersion) run(_ *cobra.Command, args []string) error {
	logger := common.GetLogger()
	schedulerName := args[0]
	targetVersion := args[0]

	request := v1.SwitchActiveVersionRequest{
		Version: targetVersion,
	}

	serializedRequest, err := common.Marshaller.Marshal(&request)
	if err != nil {
		return fmt.Errorf("error parsing request to json: %w", err)
	}

	logger.Debug("switch active version to scheduler: " + schedulerName)

	url := fmt.Sprintf("%s/schedulers/%s", a.config.ServerURL, schedulerName)
	body, status, err := a.client.Put(url, string(serializedRequest))
	if err != nil {
		return fmt.Errorf("error on put request: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("switch active version response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}
	var response v1.SwitchActiveVersionResponse
	err = protojson.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error deserializing switch active version response, details: %w", err)
	}
	logger.Info("Successfully executed switch active version operation, operation id: " + response.OperationId)
	return nil
}
