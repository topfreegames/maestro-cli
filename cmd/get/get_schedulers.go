// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package get

import (
	"fmt"
	"net/http"
	"os"
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

// getSchedulersCmd represents the list command
var getSchedulersCmd = &cobra.Command{
	Use:     "schedulers",
	Short:   "Lists all schedulers",
	Example: "maestro-cli get schedulers",
	Long:    "Lists all schedulers of a given context.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetSchedulers(client, config).run(cmd, args)
	},
}

type GetSchedulers struct {
	client interfaces.Client
	config *extensions.Config
}

func NewGetSchedulers(client interfaces.Client, config *extensions.Config) *GetSchedulers {
	return &GetSchedulers{
		client: client,
		config: config,
	}
}

func (cs *GetSchedulers) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()

	logger.Debug("getting schedulers")

	url := fmt.Sprintf("%s/schedulers", cs.config.ServerURL)
	body, status, err := cs.client.Get(url, "")
	if err != nil {
		return fmt.Errorf("error on GET request: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("get schedulers response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	logger.Sugar().Debugf("success getting schedulers: %s", body)

	var schedulers v1.ListSchedulersResponse
	err = protojson.Unmarshal(body, &schedulers)
	if err != nil {
		return fmt.Errorf("error parsing response body: %w", err)
	}

	cs.printSchedulersTable(schedulers.Schedulers)

	return nil
}

func (cs *GetSchedulers) printSchedulersTable(schedulers []*v1.Scheduler) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	format := "\n %s\t\t%s\t\t%s\t\t%s\t\t%s\t"
	fmt.Fprintf(w, format, "GAME", "NAME", "STATE", "VERSION", "AGE")

	for _, scheduler := range schedulers {
		age := time.Since(scheduler.GetCreatedAt().AsTime())
		prettyAge := durafmt.ParseShort(age).String()
		fmt.Fprintf(w, format, scheduler.GetGame(), scheduler.GetName(), scheduler.GetState(), scheduler.GetVersion(), prettyAge)
	}
}
