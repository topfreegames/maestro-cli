// maestro-cli
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

	. "github.com/topfreegames/maestro-cli/login"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {
	var (
		login     *Login
		serverURL = "server.com"
		path      string
	)

	BeforeEach(func() {
		login = NewLogin(serverURL, func(string) error {
			return nil
		})
		path = fmt.Sprintf("%s/login?state=%s", serverURL, login.OAuthState)
	})

	It("should perform login with success", func() {
		resp, err := json.Marshal(map[string]interface{}{
			"url": "google.authenticate.com",
		})
		Expect(err).NotTo(HaveOccurred())

		client.EXPECT().Get(path).Return(resp, http.StatusOK, nil)

		err = login.Perform(client)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should return error if get request fails", func() {
		client.EXPECT().Get(path).Return(nil, 0, errors.New("request error"))

		err := login.Perform(client)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("request error"))
	})

	It("should return error is status code is not 200", func() {
		client.EXPECT().Get(path).Return(nil, http.StatusBadRequest, nil)

		err := login.Perform(client)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("status code 400 when GET request to controller server"))
	})

	It("should return error if browser fails", func() {
		resp, err := json.Marshal(map[string]interface{}{
			"url": "google.authenticate.com",
		})
		Expect(err).NotTo(HaveOccurred())

		login := NewLogin(serverURL, func(string) error {
			return errors.New("browser error")
		})
		path = fmt.Sprintf("%s/login?state=%s", serverURL, login.OAuthState)

		client.EXPECT().Get(path).Return(resp, http.StatusOK, nil)

		err = login.Perform(client)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("browser error"))
	})
})
