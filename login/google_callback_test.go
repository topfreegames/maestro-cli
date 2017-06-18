// maestro-cli api
// +build unit
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package login_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/afero"

	. "github.com/topfreegames/maestro-cli/login"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoogleCallback", func() {
	var (
		state         = "state"
		expectedState = "state"
		code          = "code"
		serverURL     = "server.example.com"
		token         = "i-am-a-token"
		configDir     string
		configPath    string
		url           string
		context       = "test"
	)

	BeforeEach(func() {
		usr, err := user.Current()
		Expect(err).NotTo(HaveOccurred())
		configDir = filepath.Join(usr.HomeDir, ".maestro")
		configPath = filepath.Join(configDir, "config-test.yaml")
		url = fmt.Sprintf("%s/access?code=%s", serverURL, code)
	})

	It("should get access token and save it locally", func() {
		resp, err := json.Marshal(map[string]interface{}{
			"token": token,
		})
		Expect(err).NotTo(HaveOccurred())

		client.EXPECT().Get(url).
			Return(resp, http.StatusOK, nil)
		filesystem.EXPECT().
			MkdirAll(configDir, os.ModePerm).
			Return(nil)

		file, err := afero.NewMemMapFs().Create(configPath)
		Expect(err).NotTo(HaveOccurred())

		filesystem.EXPECT().
			Create(configPath).
			Return(file, nil)

		err = SaveAccessToken(
			state, code, expectedState, serverURL,
			context,
			filesystem,
			client,
		)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return error if state is not expectedState", func() {
		state := "random-state"
		err := SaveAccessToken(
			state, code, expectedState, serverURL,
			context,
			filesystem,
			client,
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("invalid oauth state, expected 'state', got 'random-state'"))
	})

	It("should return error if get request fails", func() {
		client.EXPECT().Get(url).
			Return(nil, 0, errors.New("request error"))

		err := SaveAccessToken(
			state, code, expectedState, serverURL,
			context,
			filesystem,
			client,
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("request error"))
	})

	It("should return error if get request not returns status code 200", func() {
		client.EXPECT().Get(url).
			Return([]byte("bad request"), http.StatusBadRequest, nil)

		err := SaveAccessToken(
			state, code, expectedState, serverURL,
			context,
			filesystem,
			client,
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("bad request"))
	})
})
