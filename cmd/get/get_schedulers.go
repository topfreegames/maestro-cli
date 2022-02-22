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

var getSchedulersName, getSchedulersGame, getSchedulersVersion string

// getSchedulersCmd represents the list command
var getSchedulersCmd = &cobra.Command{
	Use:     "schedulers",
	Short:   "Lists all schedulers",
	Example: "maestro-cli get schedulers",
	Long:    "Lists all schedulers of a given context.",
	RunE: func(cmd *cobra.Command, args []string) error {
		parameters := &GetSchedulersParameters{
			Name:    getSchedulersName,
			Game:    getSchedulersGame,
			Version: getSchedulersVersion,
		}

		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetSchedulers(client, config, parameters).run(cmd, args)
	},
}

type GetSchedulersParameters struct {
	Name    string
	Game    string
	Version string
}

type GetSchedulers struct {
	client     interfaces.Client
	config     *extensions.Config
	parameters *GetSchedulersParameters
}

func init() {
	getSchedulersCmd.Flags().StringVarP(&getSchedulersName, "name", "n", "", "Add name filter")
	getSchedulersCmd.Flags().StringVarP(&getSchedulersGame, "game", "g", "", "Add game filter")
	getSchedulersCmd.Flags().StringVarP(&getSchedulersVersion, "version", "t", "", "Add version filter")
}

func NewGetSchedulers(client interfaces.Client, config *extensions.Config, parameters *GetSchedulersParameters) *GetSchedulers {
	return &GetSchedulers{
		client:     client,
		config:     config,
		parameters: parameters,
	}
}

func (cs *GetSchedulers) run(_ *cobra.Command, args []string) error {

	logger := common.GetLogger()

	logger.Debug("getting schedulers")

	parameters := buildURLParameters(cs.parameters.Name, cs.parameters.Game, cs.parameters.Version)

	url := fmt.Sprintf("%s/schedulers%s", cs.config.ServerURL, parameters)
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

func (cs *GetSchedulers) printSchedulersTable(schedulers []*v1.SchedulerWithoutSpec) {
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

func buildURLParameters(name, game, version string) string {
	parameters := ""

	parameters = appendParameter(parameters, "name", name)
	parameters = appendParameter(parameters, "game", game)
	parameters = appendParameter(parameters, "version", version)

	return parameters
}

func appendParameter(parameters, parameterName, parameter string) string {
	if parameter != "" {
		if len(parameters) == 0 {
			parameters = "?"
		} else {
			parameters = parameters + "&"
		}
		parameters = parameters + fmt.Sprintf("%s=%s", parameterName, parameter)
	}

	return parameters
}
