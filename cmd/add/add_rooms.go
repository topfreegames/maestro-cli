// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package add

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"
)

// addRoomsCmd represents the create command
var addRoomsCmd = &cobra.Command{
	Use:     "rooms",
	Short:   "Add rooms to a given scheduler",
	Example: "maestro-cli add rooms <scheduler_name> <amount>",
	Long:    "Given the scheduler name and the amount, increase the number of rooms in Maestro.",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewAddRooms(client, config).run(cmd, args)
	},
}

type AddRooms struct {
	client interfaces.Client
	config *extensions.Config
}

func NewAddRooms(client interfaces.Client, config *extensions.Config) *AddRooms {
	return &AddRooms{
		client: client,
		config: config,
	}
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

func (a *AddRooms) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()
	schedulerName := args[0]
	roomsAmount, _ := strconv.ParseInt(args[1], 10, 32)

	request := v1.AddRoomsRequest{
		Amount: int32(roomsAmount),
	}

	serializedRequest, err := common.Marshaller.Marshal(&request)
	if err != nil {
		return fmt.Errorf("error parsing request to json: %w", err)
	}

	logger.Debug("addding rooms to scheduler: " + schedulerName)

	url := fmt.Sprintf("%s/schedulers/%s/add-rooms", a.config.ServerURL, schedulerName)
	body, status, err := a.client.Post(url, string(serializedRequest))
	if err != nil {
		return fmt.Errorf("error on post request: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("add rooms response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	logger.Info("Successfully rooms added: " + schedulerName)

	return nil
}
