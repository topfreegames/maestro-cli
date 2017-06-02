// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package interfaces

//Client interface
type Client interface {
	Get(url string) ([]byte, int, error)
	Put(url string, body map[string]interface{}) ([]byte, int, error)
	Post(url string, body map[string]interface{}) ([]byte, int, error)
	Delete(url string) ([]byte, int, error)
}
