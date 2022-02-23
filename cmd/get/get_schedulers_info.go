package get

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"github.com/topfreegames/maestro-cli/extensions"
	"github.com/topfreegames/maestro-cli/interfaces"
	v1 "github.com/topfreegames/maestro/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var getSchedulersInfoCmd = &cobra.Command{
	Use:     "schedulers-info",
	Short:   "List information from schedulers and game rooms",
	Example: "maestro-cli get schedulers-info <game>",
	Long:    "Lists schedulers and game rooms information for all schedulers or from specific game.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, config, err := common.GetClientAndConfig()
		if err != nil {
			return err
		}

		return NewGetSchedulersInfo(client, config).run(cmd, args)
	},
}

type GetSchedulersInfo struct {
	client interfaces.Client
	config *extensions.Config
}

func NewGetSchedulersInfo(client interfaces.Client, config *extensions.Config) *GetSchedulersInfo {
	return &GetSchedulersInfo{
		client: client,
		config: config,
	}
}

func (s *GetSchedulersInfo) run(_ *cobra.Command, args []string) error {
	logger := common.GetLogger()
	var game string
	if len(args) > 0 {
		game = args[0]
	}

	var url string
	if game != "" {
		logger.Debug("get schedulers information to game: " + game)
		url = fmt.Sprintf("%s/schedulers/info?game=/%s", s.config.ServerURL, game)
	} else {
		logger.Debug("get schedulers information to all schedulers")
		url = fmt.Sprintf("%s/schedulers/info", s.config.ServerURL)
	}
	body, status, err := s.client.Get(url, "")
	if err != nil {
		return fmt.Errorf("error on GET request: %w", err)
	}
	if status != http.StatusOK {
		return fmt.Errorf("get schedulers info by game response not ok, status: %s, body: %s", http.StatusText(status), string(body))
	}

	var schedulers v1.GetSchedulersInfoResponse
	err = protojson.Unmarshal(body, &schedulers)
	if err != nil {
		return fmt.Errorf("error parsing response body: %w", err)
	}

	s.printSchedulersTable(schedulers.Schedulers)

	return nil
}

func (s *GetSchedulersInfo) printSchedulersTable(schedulers []*v1.SchedulerInfo) {
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()

	format := "\n %s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%s\t"
	fmt.Fprintf(w, format, "SCHEDULER", "GAME", "STATE", "ROOM_SREADY", "ROOMS_OCCUPIED", "ROOMS_CREATING", "ROOMS_TERMINATING")

	for _, scheduler := range schedulers {
		fmt.Fprintf(w, format, scheduler.GetName(), scheduler.GetName(), scheduler.GetState(), strconv.Itoa(int(scheduler.GetRoomsReady())), strconv.Itoa(int(scheduler.GetRoomsOccupied())), strconv.Itoa(int(scheduler.GetRoomsCreating())), strconv.Itoa(int(scheduler.GetRoomsTerminating())))
	}
}
