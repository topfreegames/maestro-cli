// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions

import (
	"github.com/spf13/afero"
	"os"
)

type FileSystem struct{}

func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

//MkdirAll creates a mock directory
func (m *FileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

//Create creates a mock file
func (m *FileSystem) Create(name string) (afero.File, error) {
	return os.Create(name)
}

//IsNotExist returns true if err if of type FileNotExists
func (m *FileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

//Stat returns the FileInfo describing the the file
func (m *FileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}
