package get

import (
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/mocks"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"testing"
)

func TestGetOperationAction(t *testing.T) {
	type input struct {
		schedulerName      string
		operationId        string
		includeInputFlag   bool
		includeHistoryFlag bool
	}
	type mockPreparation struct {
		operationRetrieved *v1.GetOperationResponse
		statusCode         int
		clientError        error
	}
	tests := []struct {
		name            string
		input           input
		mockPreparation mockPreparation
		errWanted       error
	}{
		{
			name: "return no error when operation is retrieved with success using flags",
			input: input{
				schedulerName:      "scheduler_name",
				operationId:        "operation_id",
				includeInputFlag:   true,
				includeHistoryFlag: true,
			},
			mockPreparation: mockPreparation{
				operationRetrieved: &v1.GetOperationResponse{
					Operation: &v1.Operation{
						Id:             "123",
						Status:         "finished",
						DefinitionName: "remove_rooms",
						Lease:          &v1.Lease{Ttl: "2022-01-01"},
						SchedulerName:  "Name",
						CreatedAt:      timestamppb.Now(),
						Input: &_struct.Struct{
							Fields: map[string]*structpb.Value{
								"value": {
									Kind: &structpb.Value_StringValue{
										StringValue: "c50acc91-4d88-46fa-aa56-48d63c5b5311",
									},
								},
							},
						},
						ExecutionHistory: []*v1.OperationEvent{
							{
								CreatedAt: timestamppb.Now(),
								Event:     "started",
							},
						},
					},
				},
				statusCode:  200,
				clientError: nil,
			},
			errWanted: nil,
		},
		{
			name: "return no error when operation is retrieved with success not using flags",
			input: input{
				schedulerName:      "scheduler_name",
				operationId:        "operation_id",
				includeInputFlag:   false,
				includeHistoryFlag: false,
			},
			mockPreparation: mockPreparation{
				operationRetrieved: &v1.GetOperationResponse{
					Operation: &v1.Operation{
						Id:             "123",
						Status:         "finished",
						DefinitionName: "remove_rooms",
						Lease:          &v1.Lease{Ttl: "2022-01-01"},
						SchedulerName:  "Name",
						CreatedAt:      timestamppb.Now(),
						Input: &_struct.Struct{
							Fields: map[string]*structpb.Value{
								"value": {
									Kind: &structpb.Value_StringValue{
										StringValue: "c50acc91-4d88-46fa-aa56-48d63c5b5311",
									},
								},
							},
						},
						ExecutionHistory: []*v1.OperationEvent{
							{
								CreatedAt: timestamppb.Now(),
								Event:     "started",
							},
						},
					},
				},
				statusCode:  200,
				clientError: nil,
			},
			errWanted: nil,
		},

		{
			name: "return error when get request fails",
			input: input{
				schedulerName: "scheduler_name",
				operationId:   "operation_id",
			},
			mockPreparation: mockPreparation{
				operationRetrieved: nil,
				statusCode:         500,
				clientError:        errors.New("some client error"),
			},
			errWanted: errors.New("error on GET request: some client error"),
		},
		{
			name: "return error when status code is not 200",
			input: input{
				schedulerName: "scheduler_name",
				operationId:   "operation_id",
			},
			mockPreparation: mockPreparation{
				operationRetrieved: nil,
				statusCode:         400,
				clientError:        nil,
			},
			errWanted: errors.New("get operation response not ok, status: Bad Request, body: {}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			config := &extensions.Config{ServerURL: "http://localhost:8080"}
			client := mocks.NewMockClient(mockCtrl)
			responseBody, err := protojson.Marshal(tt.mockPreparation.operationRetrieved)
			require.NoError(t, err)
			client.EXPECT().
				Get(fmt.Sprintf("%s/schedulers/%s/operations/%s", config.ServerURL, tt.input.schedulerName, tt.input.operationId), gomock.Any()).Return(responseBody, tt.mockPreparation.statusCode, tt.mockPreparation.clientError)
			includeOperationInput = tt.input.includeInputFlag
			includeOperationExecutionHistory = tt.input.includeHistoryFlag

			err = NewGetOperation(client, config).runGetOperation(&cobra.Command{}, []string{tt.input.schedulerName, tt.input.operationId})

			if tt.errWanted != nil {
				require.EqualError(t, err, tt.errWanted.Error())
			} else {
				require.NoError(t, err)
			}

		})
	}
}
