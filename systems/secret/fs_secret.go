package secret

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-home-io/server/plugins/common"
	"github.com/go-home-io/server/plugins/secret"
	"github.com/go-home-io/server/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type fsSecret struct {
	location string
	logger   common.ILoggerProvider
	secrets  map[string]string
}

// Init loads initial config file.
func (s *fsSecret) Init(data *secret.InitDataSecret) error {
	s.logger = data.Logger
	loc, ok := data.Options["location"]
	if !ok {
		loc = fmt.Sprintf("%s/configs/_secrets.yaml", utils.GetCurrentWorkingDir())
		s.logger.Info("Using default location", "location", loc)
	}

	s.location = loc
	s.secrets = make(map[string]string)

	fileData, err := ioutil.ReadFile(s.location)
	if err != nil {
		s.logger.Error("Failed to read _secrets.yaml file", err)
		return nil
	}

	err = yaml.Unmarshal(fileData, &s.secrets)
	if err != nil {
		s.logger.Error("Failed to unmarshal _secrets.yaml file. Be aware that it will be rewritten", err)
		s.secrets = make(map[string]string)
		return nil
	}

	return nil
}

// Get returns found secret or throws an error.
func (s *fsSecret) Get(name string) (string, error) {
	key, ok := s.secrets[name]
	if !ok {
		return "", errors.New("not found")
	}

	return key, nil
}

// Set adds a new secret and overwrites existing file.
func (s *fsSecret) Set(name string, data string) error {
	s.secrets[name] = data
	fileData, err := yaml.Marshal(s.secrets)
	if err != nil {
		s.logger.Error("Failed to marshal secrets data", err)
		return err
	}

	writer, err := os.Create(s.location)
	if err != nil {
		s.logger.Warn("Failed to open _secrets.yaml file")
		return err
	}

	defer writer.Close()
	_, err = writer.Write(fileData)
	if err != nil {
		s.logger.Warn("Failed to write _secrets.yaml file")
		return err
	}

	return nil
}

// UpdateLogger required for logger update, since this provider is loaded before the main logger.
func (s *fsSecret) UpdateLogger(provider common.ILoggerProvider) {
	s.logger = provider
}
