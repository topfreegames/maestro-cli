// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package create

import (
	// "context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"github.com/topfreegames/maestro-cli/common"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/topfreegames/maestro/pkg/api/v1"
	yaml "gopkg.in/yaml.v2"
)

var marshler = &runtime.HTTPBodyMarshaler{
	Marshaler: &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: true,
		},
	},
}

// createSchedulerCmd represents the create command
var createSchedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Creates new scheduler",
	Long: `Creates a new scheduler on Maestro and, if worker is running, the 
	rooms will be launched.`,
	Args: validateArgs,
	RunE: run,
}

func validateArgs(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("error: missing arg with scheduler config file path")
	}

	filePath := args[0]
	_, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	return nil
}

func run(_ *cobra.Command, args []string) error {

	filePath := args[0]
	bts, _ := ioutil.ReadFile(filePath)
	
	requests := []*v1.CreateSchedulerRequest{}
	err := yaml.Unmarshal(bts, bts)
	if err != nil {
		return fmt.Errorf("error parsing file to request: %w", err)
	}

	for _, request := range requests {

		serializedRequest, err := marshler.Marshal(request)
		if err != nil {
			return fmt.Errorf("error parsing request to json: %w", err)
		}
	
		config, err := common.GetConfig()
		if err != nil {
			return fmt.Errorf("error getting client config: %w", err)
		}
	
		client := common.GetClient(config)
	
		fmt.Println("creating scheduler...")
	
		url := fmt.Sprintf("%s/schedulers", config.ServerURL)
		body, status, err := client.Post(url, string(serializedRequest))
		if err != nil {
			return fmt.Errorf("error on post request: %w", err)
		}
		if status != http.StatusOK {
			return fmt.Errorf("create scheduler response not ok, status: %s, body: %s", http.StatusText(status), string(body))
		}
	
		fmt.Println("Successfully created scheduler")
		fmt.Println(string(bts))
	}

	return nil
}
