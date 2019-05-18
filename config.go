package sptty

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	confPath string
	cfg      SpttyConfig
}

func (s *Config) Init(app Sptty) error {

	f, err := os.Open(s.confPath)
	defer f.Close()

	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &s.cfg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Config) Release() {

}

func (s *Config) AddConfigs(cfgs SpttyConfig) {
	s.cfg = cfgs
}

func (s *Config) SetConf(conf string) {
	s.confPath = conf
}

func (s *Config) Config() interface{} {
	return s.Config()
}

func (s *Config) ServiceName() string {
	return "config"
}
