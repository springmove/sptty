package sptty

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
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

type BaseConfig struct {
	Config
}

func (s *BaseConfig) ConfigName() string {
	return ""
}

func (s *BaseConfig) Validate() error {
	return nil
}

func (s *BaseConfig) Default() interface{} {
	return &BaseConfig{}
}
