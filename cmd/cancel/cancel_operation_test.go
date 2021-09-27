// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package cancel

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
)

func TestCancelOperationAction(t *testing.T) {

	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("fails when not enough args", func(t *testing.T) {
		err := validateArgs(nil, []string{})

		require.Error(t, err)
		require.Equal(t, "missing args: scheduler name or/and operation ID", err.Error())

		err = validateArgs(nil, []string{"scheduler"})

		require.Error(t, err)
		require.Equal(t, "missing args: scheduler name or/and operation ID", err.Error())
	})

	t.Run("validate args with success", func(t *testing.T) {
		err := validateArgs(nil, []string{"scheduler", uuid.NewString()})

		require.NoError(t, err)
	})

	t.Run("with success", func(t *testing.T) {

		schedulerName := "scheduler-name-1"
		operationID := "operation-id-1"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		url := fmt.Sprintf("%s/schedulers/%s/operations/%s/cancel", config.ServerURL, schedulerName, operationID)
		client.EXPECT().Post(url, gomock.Any()).Return([]byte("{}"), 200, nil)

		err := NewCancelOperation(client, config).run(nil, []string{schedulerName, operationID})

		require.NoError(t, err)
	})

	t.Run("fails when response deserialization fails", func(t *testing.T) {

		schedulerName := "scheduler-name-1"
		operationID := "operation-id-1"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		url := fmt.Sprintf("%s/schedulers/%s/operations/%s/cancel", config.ServerURL, schedulerName, operationID)
		client.EXPECT().Post(url, gomock.Any()).Return([]byte(""), 200, nil)

		err := NewCancelOperation(client, config).run(nil, []string{schedulerName, operationID})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error deserializing cancel operation response")
	})

	t.Run("fails when HTTP request fails", func(t *testing.T) {

		schedulerName := "scheduler-name-1"
		operationID := "operation-id-1"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		url := fmt.Sprintf("%s/schedulers/%s/operations/%s/cancel", config.ServerURL, schedulerName, operationID)
		client.EXPECT().Post(url, gomock.Any()).Return([]byte(""), 0, fmt.Errorf("tcp connection failed"))

		err := NewCancelOperation(client, config).run(nil, []string{schedulerName, operationID})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on post request: tcp connection failed")
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {

		schedulerName := "scheduler-name-1"
		operationID := "operation-id-1"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		url := fmt.Sprintf("%s/schedulers/%s/operations/%s/cancel", config.ServerURL, schedulerName, operationID)
		client.EXPECT().Post(url, gomock.Any()).Return([]byte(""), 400, nil)

		err := NewCancelOperation(client, config).run(nil, []string{schedulerName, operationID})

		require.Error(t, err)
		require.Contains(t, err.Error(), "cancel operation response not ok, status: Bad Request, body: ")
	})
}
