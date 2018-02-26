// Copyright © 2017 TopFreeGames backend@tfgco.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
)

func newLog(cmdName string) *logrus.Logger {
	ll := logrus.InfoLevel
	switch Verbose {
	case 0:
		ll = logrus.ErrorLevel
		break
	case 1:
		ll = logrus.WarnLevel
		break
	case 3:
		ll = logrus.DebugLevel
		break
	default:
		ll = logrus.InfoLevel
	}

	log := logrus.New()
	log.Level = ll

	return log
}

func printError(bodyResp []byte) {
	bodyMap := make(map[string]interface{})
	err := json.Unmarshal(bodyResp, &bodyMap)
	if err != nil {
		fmt.Println(string(bodyResp))
		return
	}

	errMsg, contains := bodyMap["error"]
	msg := errMsg.(string)
	if contains && strings.Contains(msg, "access token") {
		fmt.Println("You are not logged in. Please log in and try again.")
		return
	}

	for key, value := range bodyMap {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func buildBar(title string) string {
	bar := ""
	for i := 0; i < len(title); i++ {
		bar = bar + "="
	}
	return bar
}
