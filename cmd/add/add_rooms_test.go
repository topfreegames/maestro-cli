// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package add

import (
	"fmt"
	"testing"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"

	"github.com/topfreegames/maestro-cli/common"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
)

func TestAddRoomsAction(t *testing.T) {

	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("fails when not enough args", func(t *testing.T) {
		err := validateArgs(nil, []string{})

		require.Error(t, err)
		require.Equal(t, "missing args: scheduler name or/and rooms amount", err.Error())

		err = validateArgs(nil, []string{"scheduler"})

		require.Error(t, err)
		require.Equal(t, "missing args: scheduler name or/and rooms amount", err.Error())
	})

	t.Run("fails when room amount arg is not an integer", func(t *testing.T) {
		err := validateArgs(nil, []string{"scheduler", "test"})

		require.Error(t, err)
		require.Equal(t, "rooms amount must be and integer value (32 bits)", err.Error())
	})

	t.Run("validate args with success", func(t *testing.T) {
		err := validateArgs(nil, []string{"scheduler", "10"})

		require.NoError(t, err)
	})

	t.Run("with success", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		request := v1.AddRoomsRequest{
			Amount: int32(10),
		}

		serializedRequest, err := common.Marshaller.Marshal(&request)
		require.NoError(t, err)

		client.EXPECT().Post(config.ServerURL+"/schedulers/scheduler/add-rooms", string(serializedRequest)).Return([]byte(""), 200, nil)

		err = NewAddRooms(client, config).run(nil, []string{"scheduler", "10"})

		require.NoError(t, err)
	})

	t.Run("fails when HTTP request fails", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		request := v1.AddRoomsRequest{
			Amount: int32(10),
		}

		serializedRequest, err := common.Marshaller.Marshal(&request)
		require.NoError(t, err)

		client.EXPECT().Post(config.ServerURL+"/schedulers/scheduler/add-rooms", string(serializedRequest)).Return([]byte(""), 0, fmt.Errorf("tcp connection failed"))

		err = NewAddRooms(client, config).run(nil, []string{"scheduler", "10"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on post request: tcp connection failed")
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		request := v1.AddRoomsRequest{
			Amount: int32(10),
		}

		serializedRequest, err := common.Marshaller.Marshal(&request)
		require.NoError(t, err)

		client.EXPECT().Post(config.ServerURL+"/schedulers/scheduler/add-rooms", string(serializedRequest)).Return([]byte(""), 404, nil)

		err = NewAddRooms(client, config).run(nil, []string{"scheduler", "10"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "add rooms response not ok, status: Not Found, body: ")
	})
}
