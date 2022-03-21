// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package get

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/hako/durafmt"
	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var getOperationsInput, getOperationsExecutionHistory bool

// getOperationsCmd represents the list command
var getOperationsCmd = &cobra.Command{
	Use:     "operations",
	Short:   "Lists all operations from a scheduler",
	Example: "maestro-cli get operations SCHEDULER_NAME",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		parameters := &GetOperationsParameters{
			Input:            getOperationsInput,
			ExecutionHistory: getOperationsExecutionHistory,
		}

		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetOperations(client, config, parameters).run(cmd, args)
	},
}

func init() {
	getOperationsCmd.Flags().BoolVarP(&getOperationsInput, "input", "i", false, "shows input for operations")
	getOperationsCmd.Flags().BoolVarP(&getOperationsExecutionHistory, "execution-history", "x", false, "shows execution history for operations")
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("missing arg: scheduler name")
	}

	return nil
}

type GetOperationsParameters struct {
	Input            bool
	ExecutionHistory bool
}

type GetOperations struct {
	client     interfaces.Client
	config     *extensions.Config
	parameters *GetOperationsParameters
}

func NewGetOperations(client interfaces.Client, config *extensions.Config, parameters *GetOperationsParameters) *GetOperations {
	return &GetOperations{
		client:     client,
		config:     config,
		parameters: parameters,
	}
}

func (cs *GetOperations) run(_ *cobra.Command, args []string) error {
	logger := common.GetLogger()
	logger.Debug("getting operations")

	schedulerName := args[0]
	url := fmt.Sprintf("%s/schedulers/%s/operations", cs.config.ServerURL, schedulerName)
	body, status, err := cs.client.Get(url, "")
	if err != nil {
		return fmt.Errorf("error on GET request: %w", err)
	}

	if status != http.StatusOK {
		return fmt.Errorf("get operations response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	logger.Sugar().Debugf("success getting scheduler operations: %s", body)

	var operationsLists v1.ListOperationsResponse
	err = protojson.Unmarshal(body, &operationsLists)
	if err != nil {
		return fmt.Errorf("error parsing response body: %w", err)
	}

	// merge all operations into a single slice
	// TODO(gabriel.corado): add option to only show operations with specific
	// status.
	mergedOperations := append(operationsLists.GetPendingOperations(), operationsLists.GetActiveOperations()...)
	mergedOperations = append(mergedOperations, operationsLists.GetFinishedOperations()...)

	// TODO(gabriel.corado): add a option to reverse this order.
	sort.Slice(mergedOperations, func(i, j int) bool {
		return mergedOperations[i].GetCreatedAt().AsTime().Before(mergedOperations[j].GetCreatedAt().AsTime())
	})

	if len(mergedOperations) == 0 {
		fmt.Println("no operations found")
		return nil
	}

	cs.printOperationsTable(mergedOperations)
	return nil
}

// printOperationsTable receives the operations sorted and print the table to
// stdout.
func (cs *GetOperations) printOperationsTable(operations []*v1.Operation) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	defaultHeaders := []interface{}{"ID", "NAME", "STATUS", "AGE", "LEASE_TTL", "LEASE_EXPIRED"}
	format := "%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s"

	headers := make([]interface{}, 0)
	headers = append(headers, defaultHeaders...)

	if cs.parameters.Input {
		format = format + "\t\t%s"
		headers = append(headers, "INPUT")
	}
	if cs.parameters.ExecutionHistory {
		format = format + "\t\t%s"
		headers = append(headers, "EXEC. HIST.")
	}

	format = format + "\t\n"
	fmt.Fprintf(w, format, headers...)

	for _, operation := range operations {
		values := cs.getOperationValues(operation)
		fmt.Fprintf(w, format, values...)
	}
}

func (cs *GetOperations) getOperationValues(operation *v1.Operation) []interface{} {
	age := time.Since(operation.GetCreatedAt().AsTime())
	prettyAge := durafmt.ParseShort(age).String()

	leaseTtl, leaseExpired := cs.getOperationLeaseInfo(operation)

	values := make([]interface{}, 0)

	defaultValues := []interface{}{
		operation.GetId(),
		strings.ToUpper(operation.GetDefinitionName()),
		strings.ToUpper(operation.GetStatus()),
		prettyAge,
		leaseTtl,
		leaseExpired,
	}
	values = append(values, defaultValues...)
	values = cs.appendInputIfRequested(values, operation)
	values = cs.appendExecutionHistoryValueIfRequested(values, operation)
	return values
}

func (cs *GetOperations) appendExecutionHistoryValueIfRequested(values []interface{}, operation *v1.Operation) []interface{} {
	if cs.parameters.ExecutionHistory {
		type printableEvent struct {
			CreatedAt string `json:"createdAt"`
			Event     string `json:"event"`
		}

		history := make([]*printableEvent, 0)
		execHist := operation.GetExecutionHistory()
		for _, event := range execHist {
			history = append(history, &printableEvent{
				CreatedAt: event.GetCreatedAt().AsTime().String(),
				Event:     event.GetEvent(),
			})
		}
		values = append(values, cs.fromFieldToJson(history))
	}
	return values
}

func (cs *GetOperations) appendInputIfRequested(values []interface{}, operation *v1.Operation) []interface{} {
	if cs.parameters.Input {
		values = append(values, cs.fromFieldToJson(operation.GetInput()))
	}
	return values
}

func (cs *GetOperations) getOperationLeaseInfo(operation *v1.Operation) (string, string) {
	leaseTtl := "-"
	leaseExpired := "-"
	if operation.Lease != nil {
		leaseTtl = operation.Lease.GetTtl()
		parsedLeaseTtl, err := time.Parse(time.RFC3339, operation.Lease.GetTtl())
		if err == nil {
			expiredSeconds := time.Since(parsedLeaseTtl).Seconds()
			leaseExpired = strings.ToUpper(strconv.FormatBool(expiredSeconds > 0))
		}
	}
	return leaseTtl, leaseExpired
}

func (cs *GetOperations) fromFieldToJson(field interface{}) string {
	var prettyField string
	out, err := json.Marshal(field)
	if err != nil {
		return ""
	}
	prettyField = string(out)
	if prettyField == "null" {
		prettyField = "-"
	}

	return prettyField
}
