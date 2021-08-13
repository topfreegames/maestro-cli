// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package get

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetSchedulersAction(t *testing.T) {
	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		schedulers := &v1.ListSchedulersReply{
			Schedulers: []*v1.Scheduler{
				{
					Name:      "scheduler-test-1",
					Game:      "the-game",
					Version:   "1.0.0",
					State:     "creating",
					CreatedAt: timestamppb.Now(),
				},
			},
		}

		responseBody, _ := protojson.Marshal(schedulers)

		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return(responseBody, 200, nil)

		err := NewGetSchedulers(client, config).run(nil, []string{})

		require.NoError(t, err)
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewGetSchedulers(client, config).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get schedulers response not ok, status: Not Found")
	})

	t.Run("fails when bad format response body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 200, nil)

		err := NewGetSchedulers(client, config).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing response body")
	})

	t.Run("fails when HTTP request faile", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 0, errors.New("request failed"))

		err := NewGetSchedulers(client, config).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on GET request: request failed")
	})
}
