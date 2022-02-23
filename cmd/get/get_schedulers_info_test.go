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
)

func TestGetSchedulersInfoAction(t *testing.T) {
	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("it should return with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		client := mocks.NewMockClient(mockCtrl)
		schedulers := &v1.GetSchedulersInfoResponse{
			Schedulers: []*v1.SchedulerInfo{
				{
					Name:             "scheduler-test-1",
					Game:             "the-game",
					State:            "creating",
					RoomsReady:       10,
					RoomsTerminating: 2,
					RoomsCreating:    5,
					RoomsOccupied:    20,
				},
			},
		}
		responseBody, _ := protojson.Marshal(schedulers)
		client.EXPECT().Get(config.ServerURL+"/schedulers/info", gomock.Any()).Return(responseBody, 200, nil)

		err := NewGetSchedulersInfo(client, config).run(nil, []string{})

		require.NoError(t, err)
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/info", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewGetSchedulersInfo(client, config).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get schedulers info by game response not ok, status: Not Found")
	})

	t.Run("fails when bad format response body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/info", gomock.Any()).Return([]byte(""), 200, nil)

		err := NewGetSchedulersInfo(client, config).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing response body")
	})

	t.Run("fails when HTTP request fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers/info", gomock.Any()).Return([]byte(""), 0, errors.New("request failed"))

		err := NewGetSchedulersInfo(client, config).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on GET request: request failed")
	})
}
