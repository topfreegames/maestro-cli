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
	"time"

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
	params := &GetOperationsParameters{Input: true, ExecutionHistory: true}

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
					Id:               "123",
					Status:           "finished",
					DefinitionName:   "remove_rooms",
					Lease:            &v1.Lease{Ttl: "2022-01-01"},
					SchedulerName:    "Name",
					CreatedAt:        timestamppb.Now(),
					Input:            nil,
					ExecutionHistory: nil,
				},
			},
		}

		schedulerName := "test"
		responseBody, err := protojson.Marshal(operations)
		require.NoError(t, err)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return(responseBody, 200, nil)

		err = NewGetOperations(client, config, params).run(nil, []string{schedulerName})
		require.NoError(t, err)
	})

	t.Run("no operations found - with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		operations := &v1.ListOperationsResponse{
			PendingOperations:  []*v1.Operation{},
			ActiveOperations:   []*v1.Operation{},
			FinishedOperations: []*v1.Operation{},
		}

		schedulerName := "test"
		responseBody, err := protojson.Marshal(operations)
		require.NoError(t, err)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return(responseBody, 200, nil)

		err = NewGetOperations(client, config, params).run(nil, []string{schedulerName})
		require.NoError(t, err)
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		schedulerName := "test"
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewGetOperations(client, config, params).run(nil, []string{schedulerName})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get operations response not ok, status: Not Found")
	})

	t.Run("fails when bad format response body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		schedulerName := "test"
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return([]byte(""), 200, nil)

		err := NewGetOperations(client, config, params).run(nil, []string{schedulerName})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing response body")
	})

	t.Run("fails when HTTP request failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		schedulerName := "test"
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return([]byte(""), 0, errors.New("request failed"))

		err := NewGetOperations(client, config, params).run(nil, []string{schedulerName})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on GET request: request failed")
	})

	t.Run("Successfully consume lease information from operations", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		operations := &v1.ListOperationsResponse{
			PendingOperations: []*v1.Operation{
				{
					Id:             "ACTIVE_OPERATION",
					DefinitionName: "remove_rooms",
					Status:         "pending",
					CreatedAt:      timestamppb.Now(),
					Lease:          &v1.Lease{Ttl: time.Now().UTC().Add(time.Hour).Format(time.RFC3339)},
				},
				{
					Id:             "EXPIRED_OPERATION",
					DefinitionName: "remove_rooms",
					Status:         "pending",
					CreatedAt:      timestamppb.Now(),
					Lease:          &v1.Lease{Ttl: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)},
				},
			},
			FinishedOperations: []*v1.Operation{
				{
					Id:             "FINISHED_OPERATION",
					DefinitionName: "remove_rooms",
					Status:         "finished",
					CreatedAt:      timestamppb.Now(),
					Lease:          nil,
				},
			},
		}

		schedulerName := "test"
		responseBody, err := protojson.Marshal(operations)
		require.NoError(t, err)
		client.EXPECT().Get(config.ServerURL+"/schedulers/"+schedulerName+"/operations", gomock.Any()).Return(responseBody, 200, nil)

		err = NewGetOperations(client, config, params).run(nil, []string{schedulerName})
		require.NoError(t, err)
	})
}
