// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package create

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
)

func TestCreateSchedulerAction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mocks.NewMockClient(mockCtrl)

	dirPath, _ := os.Getwd()
	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("fails when no file found on path", func(t *testing.T) {
		err := NewCreateScheduler(client, config).run(nil, []string{"fixtures/scheduler-config-not-found.yaml"})

		require.Error(t, err)
		require.Equal(t, "error reading scheduler file: open fixtures/scheduler-config-not-found.yaml: no such file or directory", err.Error())
	})

	t.Run("fails when file found bad format", func(t *testing.T) {
		err := NewCreateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config-bad-format.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot unmarshal !!str `name` into v1.CreateSchedulerRequest")
	})

	t.Run("with success", func(t *testing.T) {
		client.EXPECT().Post(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 200, nil)

		err := NewCreateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.NoError(t, err)
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		client.EXPECT().Post(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewCreateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "create scheduler response not ok, status: Not Found, body: ")
	})
}
