// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package update

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestUpdateSchedulerAction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mocks.NewMockClient(mockCtrl)

	dirPath, _ := os.Getwd()
	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("with success", func(t *testing.T) {
		expectedStructuredBody := v1.UpdateSchedulerResponse{
			OperationId: "<OperationID>",
		}

		expectedStringBody, _ := protojson.Marshal(&expectedStructuredBody)

		client.EXPECT().Put(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(expectedStringBody), 200, nil)

		err := NewUpdateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.NoError(t, err)
	})

	t.Run("fails when no file found on path", func(t *testing.T) {
		err := NewUpdateScheduler(client, config).run(nil, []string{"fixtures/scheduler-config-not-found.yaml"})

		require.Error(t, err)
		require.Equal(t, "error reading scheduler file: open fixtures/scheduler-config-not-found.yaml: no such file or directory", err.Error())
	})

	t.Run("fails when file found bad format", func(t *testing.T) {
		err := NewUpdateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config-bad-format.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing Json to v1.UpdateSchedulerRequest")
		require.Contains(t, err.Error(), "unexpected token \"name\"")
	})

	t.Run("fails when file is not .yaml", func(t *testing.T) {
		err := NewUpdateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/file_not_yaml.json"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "file should be .yaml")
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		client.EXPECT().Put(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewUpdateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "update scheduler response not ok, status: Not Found, body: ")
	})

	t.Run("fails when got error on calling maestro API", func(t *testing.T) {
		client.EXPECT().Put(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 0, errors.New("error on API call"))

		err := NewUpdateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on Put request: ")
	})
}
