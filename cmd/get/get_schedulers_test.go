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

	t.Run("when there is no parameter it should return with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		schedulers := &v1.ListSchedulersResponse{
			Schedulers: []*v1.SchedulerWithoutSpec{
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

		parameters := &GetSchedulersParameters{}

		err := NewGetSchedulers(client, config, parameters).run(nil, []string{})

		require.NoError(t, err)
	})

	t.Run("when there are parameters it should return with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		schedulers := &v1.ListSchedulersResponse{
			Schedulers: []*v1.SchedulerWithoutSpec{
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

		client.EXPECT().Get(config.ServerURL+"/schedulers?name=some name", gomock.Any()).Return(responseBody, 200, nil)

		parameters := &GetSchedulersParameters{Name: "some name"}

		err := NewGetSchedulers(client, config, parameters).run(nil, []string{})

		require.NoError(t, err)
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 404, nil)

		parameters := &GetSchedulersParameters{}

		err := NewGetSchedulers(client, config, parameters).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get schedulers response not ok, status: Not Found")
	})

	t.Run("fails when bad format response body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 200, nil)

		parameters := &GetSchedulersParameters{}

		err := NewGetSchedulers(client, config, parameters).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing response body")
	})

	t.Run("fails when HTTP request faile", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)
		client.EXPECT().Get(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 0, errors.New("request failed"))

		parameters := &GetSchedulersParameters{}

		err := NewGetSchedulers(client, config, parameters).run(nil, []string{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on GET request: request failed")
	})
}

func TestBuildURLParameters(t *testing.T) {
	type input struct {
		Name    string
		Game    string
		Version string
	}

	testCases := []struct {
		Title          string
		Input          input
		ExpectedReturn string
	}{
		{
			Title:          "when there is only name",
			Input:          input{Name: "some-name"},
			ExpectedReturn: "?name=some-name",
		}, {
			Title:          "when there is only game",
			Input:          input{Game: "some-game"},
			ExpectedReturn: "?game=some-game",
		}, {
			Title:          "when there is only version",
			Input:          input{Version: "some-version"},
			ExpectedReturn: "?version=some-version",
		}, {
			Title:          "when there are name and game",
			Input:          input{Name: "some-name", Game: "some-game"},
			ExpectedReturn: "?name=some-name&game=some-game",
		}, {
			Title:          "when there are name and version",
			Input:          input{Name: "some-name", Version: "some-version"},
			ExpectedReturn: "?name=some-name&version=some-version",
		}, {
			Title:          "when there are game and version",
			Input:          input{Game: "some-game", Version: "some-version"},
			ExpectedReturn: "?game=some-game&version=some-version",
		}, {
			Title:          "when there are all parameters",
			Input:          input{Game: "some-game", Version: "some-version"},
			ExpectedReturn: "?game=some-game&version=some-version",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Title, func(t *testing.T) {
			returnedParameters := buildURLParameters(testCase.Input.Name, testCase.Input.Game, testCase.Input.Version)
			require.Equal(t, testCase.ExpectedReturn, returnedParameters)
		})
	}
}
