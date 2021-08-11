// maestro-cli
// +build unit
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package login_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/maestro-cli/mocks"

	"testing"
)

var (
	filesystem *mocks.MockFileSystem
	client     *mocks.MockClient
	mockCtrl   *gomock.Controller
)

func TestLogin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Login Suite")
}

var _ = BeforeEach(func() {
	mockCtrl = gomock.NewController(GinkgoT())
	filesystem = mocks.NewMockFileSystem(mockCtrl)
	client = mocks.NewMockClient(mockCtrl)
})
