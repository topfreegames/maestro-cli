// maestro-cli
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/topfreegames/maestro-cli/interfaces"
	yaml "gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	ServerURL string `yaml:"serverUrl"`
}

// NewConfig ctor
func NewConfig(token, serverURL string) *Config {
	c := &Config{
		ServerURL: serverURL,
	}
	return c
}

func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, nil
}

func getDirPath() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Join(home, ".maestro")
	return dirPath, nil
}

func getConfigPath(context string) (string, error) {
	dir, err := getDirPath()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("config-%s.yaml", context)
	configPath := filepath.Join(dir, fileName)
	return configPath, nil
}

// ReadConfig from file
func ReadConfig(fs interfaces.FileSystem, context string) (*Config, error) {
	configPath, err := getConfigPath(context)
	if err != nil {
		return nil, err
	}
	if _, err := fs.Stat(configPath); fs.IsNotExist(err) {
		return nil, err
	}
	bts, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(bts, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Write the config file to disk
func (c *Config) Write(fs interfaces.FileSystem, context string) error {
	configPath, err := getConfigPath(context)
	if err != nil {
		return err
	}
	cfg, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	cfgDir, err := getDirPath()
	if err != nil {
		return err
	}
	err = fs.MkdirAll(cfgDir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := fs.Create(configPath)
	if err != nil {
		return err
	}
	_, err = file.Write(cfg)
	if err != nil {
		return err
	}
	return nil
}
