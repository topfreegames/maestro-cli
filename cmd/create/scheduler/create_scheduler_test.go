// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package scheduler

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestCreateSchedulerAction(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	client := mocks.NewMockClient(mockCtrl)

	dirPath, _ := os.Getwd()
	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("with success", func(t *testing.T) {
		expectedStructuredBody := v1.CreateSchedulerRequest{
			Name:     "scheduler-test",
			Game:     "game-test",
			MaxSurge: "10%",
			Spec: &v1.Spec{
				TerminationGracePeriod: 15,
				Containers: []*v1.Container{
					{
						Name:            "example",
						Image:           "alpine:3.15.0",
						Command:         []string{"sh", "-c", "while true; do sleep 1; done"},
						ImagePullPolicy: "Always",
						Environment:     []*v1.ContainerEnvironment{},
						Requests: &v1.ContainerResources{
							Memory: "20Mi",
							Cpu:    "10m",
						},
						Limits: &v1.ContainerResources{
							Memory: "20Mi",
							Cpu:    "10m",
						},
						Ports: []*v1.ContainerPort{
							{
								Name:     "default",
								Protocol: "tcp",
								Port:     80,
							},
						},
					},
				},
			},
			PortRange: &v1.PortRange{
				Start: 80,
				End:   8000,
			},
		}

		expectedStringBody, _ := protojson.Marshal(&expectedStructuredBody)

		client.EXPECT().Post(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(expectedStringBody), 200, nil)

		err := NewCreateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.NoError(t, err)
	})

	t.Run("fails when no file found on path", func(t *testing.T) {
		err := NewCreateScheduler(client, config).run(nil, []string{"fixtures/scheduler-config-not-found.yaml"})

		require.Error(t, err)
		require.Equal(t, "error reading scheduler file: open fixtures/scheduler-config-not-found.yaml: no such file or directory", err.Error())
	})

	t.Run("fails when file found bad format", func(t *testing.T) {
		err := NewCreateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config-bad-format.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error parsing Json to v1.CreateSchedulerRequest")
		require.Contains(t, err.Error(), "unexpected token \"name\"")
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		client.EXPECT().Post(config.ServerURL+"/schedulers", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewCreateScheduler(client, config).run(nil, []string{dirPath + "/fixtures/scheduler-config.yaml"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "create scheduler response not ok, status: Not Found, body: ")
	})
}
