package sptty

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	ConfigServiceName = "config"
)

type ConfigService struct {
	confPath string
	cfgs     map[string]interface{}
}

func (s *ConfigService) Init(app Sptty) error {

	f, err := os.Open(s.confPath)
	defer f.Close()

	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &s.cfgs)
	if err != nil {
		return err
	}

	return nil
}

func (s *ConfigService) Release() {

}

func (s *ConfigService) Enable() bool {
	return true
}

func (s *ConfigService) SetConfPath(conf string) {
	s.confPath = conf
}

func (s *ConfigService) ServiceName() string {
	return ConfigServiceName
}
