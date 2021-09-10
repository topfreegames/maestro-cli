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

func TestGetOperationsAction(t *testing.T) {
	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		operations := &v1.ListOperationsResponse{
			PendingOperations: []*v1.Operation{
				{
					Id:             "123",
					DefinitionName: "remove_rooms",
					Status:         "pending",
					CreatedAt:      timestamppb.Now(),
				},
			},
			FinishedOperations: []*v1.Operation{
				{
					Id:             "123",
					DefinitionName: "remove_rooms",
					Status:         "finished",
					CreatedAt:      timestamppb.Now(),
				},
			},
		}

		schedulerName := "test"
		responseBody, err := protojson.Marshal(operations)
		require.NoError(t, err)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return(responseBody, 200, nil)

		err = NewGetOperations(client, config).run(nil, []string{schedulerName})
		require.NoError(t, err)
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		schedulerName := "test"
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewGetOperations(client, config).run(nil, []string{schedulerName})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get operations response not ok, status: Not Found")
	})

	t.Run("fails when bad format response body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		schedulerName := "test"
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return([]byte(""), 200, nil)

		err := NewGetOperations(client, config).run(nil, []string{schedulerName})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing response body")
	})

	t.Run("fails when HTTP request failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		schedulerName := "test"
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return([]byte(""), 0, errors.New("request failed"))

		err := NewGetOperations(client, config).run(nil, []string{schedulerName})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on GET request: request failed")
	})
}
