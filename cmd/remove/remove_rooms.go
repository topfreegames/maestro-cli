// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package remove

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"
)

// removeRoomsCmd represents the create command
var removeRoomsCmd = &cobra.Command{
	Use:     "rooms",
	Short:   "remove rooms from a given scheduler",
	Example: "maestro-cli remove rooms <scheduler_name> <amount>",
	Long:    "Given the scheduler name and the amount, remove the number of rooms in Maestro.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewRemoveRooms(client, config).run(cmd, args)
	},
}

type RemoveRooms struct {
	client interfaces.Client
	config *extensions.Config
}

func NewRemoveRooms(client interfaces.Client, config *extensions.Config) *RemoveRooms {
	return &RemoveRooms{
		client: client,
		config: config,
	}
}

func (a *RemoveRooms) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()
	schedulerName := args[0]
	roomsAmount, _ := strconv.ParseInt(args[1], 10, 32)

	request := v1.RemoveRoomsRequest{
		Amount: int32(roomsAmount),
	}

	serializedRequest, err := protojson.Marshal(&request)
	if err != nil {
		return fmt.Errorf("error parsing request to json: %w", err)
	}

	logger.Debug("removing rooms from scheduler: " + schedulerName)

	url := fmt.Sprintf("%s/schedulers/%s/remove-rooms", a.config.ServerURL, schedulerName)
	body, status, err := a.client.Post(url, string(serializedRequest))
	if err != nil {
		return fmt.Errorf("error on post request: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("remove rooms response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}
	var response v1.RemoveRoomsResponse
	err = protojson.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error deserializing remove rooms response, details: %w", err)
	}

	logger.Info("Successfully executed remove rooms, operation id: " + response.OperationId)

	return nil
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("missing args: scheduler name or/and rooms amount")
	}

	if _, err := strconv.ParseInt(args[1], 10, 32); err != nil {
		return errors.New("rooms amount must be and integer value (32 bits)")
	}

	return nil
}
