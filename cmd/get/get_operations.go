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
	"sort"
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

// getOperationsCmd represents the list command
var getOperationsCmd = &cobra.Command{
	Use:     "operations",
	Short:   "Lists all operations from a scheduler",
	Example: "maestro-cli get operations SCHEDULER_NAME",
	Args:    validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetOperations(client, config).run(cmd, args)
	},
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("missing arg: scheduler name")
	}

	return nil
}

type GetOperations struct {
	client interfaces.Client
	config *extensions.Config
}

func NewGetOperations(client interfaces.Client, config *extensions.Config) *GetOperations {
	return &GetOperations{
		client: client,
		config: config,
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

	format := "%s\t\t%s\t\t%s\t\t%s\t\n"
	fmt.Fprintf(w, format, "ID", "NAME", "STATUS", "AGE")

	for _, operation := range operations {
		age := time.Since(operation.GetCreatedAt().AsTime())
		prettyAge := durafmt.ParseShort(age).String()

		fmt.Fprintf(
			w,
			format,
			operation.GetId(),
			strings.ToUpper(operation.GetDefinitionName()),
			strings.ToUpper(operation.GetStatus()),
			prettyAge,
		)
	}
}
