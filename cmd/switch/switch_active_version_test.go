package _switch

import (
	"fmt"
	"testing"

	v1 "github.com/topfreegames/maestro/pkg/api/v1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestSwitchSchedulerVersion(t *testing.T) {

	config := &extensions.Config{
		ServerURL: "http://localhost:8080",
	}

	t.Run("fails when not enough args", func(t *testing.T) {
		err := validateArgs(nil, []string{"scheduler"})

		require.Error(t, err)
		require.Equal(t, "missing args: scheduler name or/and target version", err.Error())
	})

	t.Run("fails when version format is invalid", func(t *testing.T) {
		err := validateArgs(nil, []string{"scheduler", "invalid-version"})

		require.Error(t, err)
		require.Equal(t, "invalid version value. version must be in a semantic format", err.Error())
	})

	t.Run("validate args with success", func(t *testing.T) {
		err := validateArgs(nil, []string{"scheduler-name", "v1.0.0"})

		require.NoError(t, err)
	})

	t.Run("with success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		client := mocks.NewMockClient(mockCtrl)
		response := v1.SwitchActiveVersionResponse{
			OperationId: "<OperationId>",
		}
		expectedResponse, err := protojson.Marshal(&response)

		require.NoError(t, err)

		client.EXPECT().Put(config.ServerURL+"/schedulers/scheduler-name", gomock.Any()).Return([]byte(expectedResponse), 200, nil)

		err = NewSwitchActiveVersion(client, config).run(nil, []string{"scheduler-name", "v1.0.0"})

		require.NoError(t, err)
	})

	t.Run("fails when HTTP request fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		client.EXPECT().Put(config.ServerURL+"/schedulers/scheduler-name", gomock.Any()).Return([]byte(""), 404, fmt.Errorf("tcp connection failed"))

		err := NewSwitchActiveVersion(client, config).run(nil, []string{"scheduler-name", "v1.0.0"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "error on put request: tcp connection failed")
	})

	t.Run("fails when maestro API fails", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		client := mocks.NewMockClient(mockCtrl)

		client.EXPECT().Put(config.ServerURL+"/schedulers/scheduler-name", gomock.Any()).Return([]byte(""), 404, nil)

		err := NewSwitchActiveVersion(client, config).run(nil, []string{"scheduler-name", "v1.0.0"})

		require.Error(t, err)
		require.Contains(t, err.Error(), "switch active version response not ok, status: Not Found, body: ")
	})
}
