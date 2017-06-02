// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package interfaces

import (
	"github.com/spf13/afero"
	"os"
)

//FileSystem interface
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (afero.File, error)
	IsNotExist(err error) bool
	Stat(name string) (os.FileInfo, error)
}
