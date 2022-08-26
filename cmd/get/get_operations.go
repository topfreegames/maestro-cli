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

var getStageName string
var getPageNumber, getPerPageNumber string

func init() {
	getOperationsCmd.Flags().StringVarP(&getStageName, "stage", "s", "", "Add name filter")
	getOperationsCmd.Flags().StringVarP(&getPageNumber, "page", "p", "", "Add game filter")
	getOperationsCmd.Flags().StringVarP(&getPerPageNumber, "perpage", "P", "", "Add version filter")
}

// getOperationsCmd represents the list command
var getOperationsCmd = &cobra.Command{
	Use:     "operations",
	Short:   "Lists all operations from a scheduler",
	Example: "maestro-cli get operations SCHEDULER_NAME",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		parameters := &GetOperationsParameters{
			Stage:   getStageName,
			Page:    getPageNumber,
			PerPage: getPerPageNumber,
		}
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetOperations(client, config, parameters).run(cmd, args)
	},
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("missing arg: scheduler name")
	}

	return nil
}

type GetOperationsParameters struct {
	Stage   string
	Page    string
	PerPage string
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
	parameters := buildURLParametersForOperations(cs.parameters.Stage, cs.parameters.Page, cs.parameters.PerPage)
	url := fmt.Sprintf("%s/schedulers/%s/operations%s", cs.config.ServerURL, schedulerName, parameters)

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
	operations := operationsLists.GetOperations()

	// TODO(gabriel.corado): add a option to reverse this order.
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].GetCreatedAt().AsTime().Before(operations[j].GetCreatedAt().AsTime())
	})

	if len(operations) == 0 {
		fmt.Println("no operations found")
		return nil
	}

	// cs.printOperationsTable(operations)
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

	format = format + "\t\n"
	fmt.Fprintf(w, format, headers...)

	for _, operation := range operations {
		values := getOperationValues(operation, false, false)
		fmt.Fprintf(w, format, values...)
	}
}

func getOperationValues(operation *v1.Operation, includeInput, incluseHistory bool) []interface{} {
	age := time.Since(operation.GetCreatedAt().AsTime())
	prettyAge := durafmt.ParseShort(age).String()

	leaseTtl, leaseExpired := getOperationLeaseInfo(operation)

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
	values = appendInputIfRequested(values, operation, includeInput)
	values = appendExecutionHistoryValueIfRequested(values, operation, incluseHistory)
	return values
}

func appendExecutionHistoryValueIfRequested(values []interface{}, operation *v1.Operation, includeHistory bool) []interface{} {
	if includeHistory {
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
		values = append(values, fromFieldToJson(history))
	}
	return values
}

func appendInputIfRequested(values []interface{}, operation *v1.Operation, includeInput bool) []interface{} {
	if includeInput {
		values = append(values, fromFieldToJson(operation.GetInput()))
	}
	return values
}

func getOperationLeaseInfo(operation *v1.Operation) (string, string) {
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

func fromFieldToJson(field interface{}) string {
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
func buildURLParametersForOperations(stage string, page, perpage string) string {
	parameters := ""

	parameters = appendParameter(parameters, "stage", stage)
	parameters = appendParameter(parameters, "page", page)
	parameters = appendParameter(parameters, "perPage", perpage)

	return parameters
}
