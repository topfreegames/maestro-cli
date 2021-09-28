// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cancel

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"
)

// cancelOperationCmd represents the cancel operation command
var cancelOperationCmd = &cobra.Command{
	Use:     "operation",
	Short:   "cancels operation from a given scheduler and ID",
	Example: "maestro-cli cancel operation <scheduler_name> <id>",
	Long:    "Given the scheduler name and the operation ID, cancels the found operation in Maestro.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewCancelOperation(client, config).run(cmd, args)
	},
}

type CancelOperation struct {
	client interfaces.Client
	config *extensions.Config
}

func NewCancelOperation(client interfaces.Client, config *extensions.Config) *CancelOperation {
	return &CancelOperation{
		client: client,
		config: config,
	}
}

func (a *CancelOperation) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()
	schedulerName := args[0]
	operationID := args[1]

	request := v1.CancelOperationRequest{}
	serializedRequest, err := protojson.Marshal(&request)
	if err != nil {
		return fmt.Errorf("error parsing request to json: %w", err)
	}

	logger.Sugar().Debugf("cancel operation %s from scheduler %s", schedulerName, operationID)

	url := fmt.Sprintf("%s/schedulers/%s/operations/%s/cancel", a.config.ServerURL, schedulerName, operationID)
	body, status, err := a.client.Post(url, string(serializedRequest))
	if err != nil {
		return fmt.Errorf("error on post request: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("cancel operation response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}
	var response v1.CancelOperationResponse
	err = protojson.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error deserializing cancel operation response, details: %w", err)
	}

	logger.Info("cancel operation request successfully sent")

	return nil
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("missing args: scheduler name or/and operation ID")
	}

	return nil
}
