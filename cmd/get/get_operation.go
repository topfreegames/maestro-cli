// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package get

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var includeOperationInput, includeOperationExecutionHistory bool

func init() {
	getOperationCmd.Flags().BoolVarP(&includeOperationInput, "input", "i", false, "shows input for operations")
	getOperationCmd.Flags().BoolVarP(&includeOperationExecutionHistory, "execution-history", "x", false, "shows execution history for operations")
}

// getOperationCmd represents the getOperation command
var getOperationCmd = &cobra.Command{
	Use:     "operation",
	Short:   "Get the specific operation info",
	Example: "maestro-cli get operation SCHEDULER_NAME OPERATION_ID",
	Args:    validateGetOperationArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetOperation(client, config).runGetOperation(cmd, args)
	},
}

func validateGetOperationArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("missing args: scheduler name and operation id")
	}
	if len(args) < 2 {
		return errors.New("missing arg: operation id")
	}

	return nil
}

type GetOperation struct {
	client interfaces.Client
	config *extensions.Config
}

func NewGetOperation(client interfaces.Client, config *extensions.Config) *GetOperation {
	return &GetOperation{
		client: client,
		config: config,
	}
}

func (cs *GetOperation) runGetOperation(_ *cobra.Command, args []string) error {
	logger := common.GetLogger()
	logger.Debug("getting operation")

	schedulerName := args[0]
	operationID := args[1]
	url := fmt.Sprintf("%s/schedulers/%s/operations/%s", cs.config.ServerURL, schedulerName, operationID)
	body, status, err := cs.client.Get(url, "")
	if err != nil {
		return fmt.Errorf("error on GET request: %w", err)
	}

	if status != http.StatusOK {
		return fmt.Errorf("get operation response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	logger.Sugar().Debugf("success getting operation %s: %s", body, operationID)

	var operationResponse v1.GetOperationResponse
	err = protojson.Unmarshal(body, &operationResponse)
	if err != nil {
		return fmt.Errorf("error parsing response body: %w", err)
	}

	cs.printOperation(operationResponse.Operation)
	return nil
}

func (cs *GetOperation) printOperation(operation *v1.Operation) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	defaultHeaders := []interface{}{"ID", "NAME", "STATUS", "AGE", "LEASE_TTL", "LEASE_EXPIRED"}
	format := "%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s"

	headers := make([]interface{}, 0)
	headers = append(headers, defaultHeaders...)

	if includeOperationInput {
		format = format + "\t\t%s"
		headers = append(headers, "INPUT")
	}
	if includeOperationExecutionHistory {
		format = format + "\t\t%s"
		headers = append(headers, "EXEC. HIST.")
	}

	format = format + "\t\n"
	fmt.Fprintf(w, format, headers...)

	values := getOperationValues(operation, includeOperationInput, includeOperationExecutionHistory)
	fmt.Fprintf(w, format, values...)
}
