package sptty

import (
	"errors"
	"flag"
	"fmt"

	"gopkg.in/yaml.v2"
)

var appService *AppService = nil
var appTag string = ""

func SetTag(tag string) {
	appTag = tag
}

func GetApp() *AppService {
	if appService == nil {
		appService = &AppService{
			config:   &ConfigService{},
			log:      &LogService{},
			services: Services{},
			configs: map[string]IConfig{
				LogServiceName: &LogConfig{},
			},
		}
	}

	return appService
}

func Log(level LogLevel, msg string, tags ...string) {
	app := GetApp()
	app.log.Log(level, msg, tags...)
}

type AppService struct {
	services Services
	configs  map[string]IConfig
	config   *ConfigService
	log      *LogService
}

func (s *AppService) init() error {
	if err := s.config.Init(s); err != nil {
		return err
	}

	if err := s.validateConfigs(); err != nil {
		return err
	}

	if err := s.log.Init(s); err != nil {
		return err
	}

	for _, v := range s.services {
		Log(InfoLevel, fmt.Sprintf("Init Service: %s", v.ServiceName()), v.ServiceName())
		if err := v.Init(s); err != nil {
			Log(ErrorLevel, fmt.Sprintf("Init Service %s failed: %s", v.ServiceName(), err.Error()), v.ServiceName())
			return err
		}
	}

	return nil
}

func (s *AppService) release() {
	for _, v := range s.services {
		v.Release()
	}

}

func (s *AppService) Sptting() {
	if err := s.init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	s.release()
}

// func (s *AppService) AddServices(services Services) {
// 	s.services = services
// }

// func (s *AppService) AddConfigs(configs Configs) {
// 	for k, v := range configs {
// 		s.configs[v.ConfigName()] = configs[k]
// 	}
// }

func (s *AppService) validateConfigs() error {
	for _, v := range s.configs {
		err := v.Validate()
		if err != nil {
			fmt.Printf("Config Error: %s\n", err.Error())
			return err
		}
	}

	return nil
}

func (s *AppService) ConfFromFile(conf string) {
	s.config.SetConfPath(conf)
}

func (s *AppService) LoadConfFromFile() {
	cfg := flag.String("config", "./config.yml", "--config")
	flag.Parse()
	s.ConfFromFile(*cfg)
}

func (s *AppService) GetConfig(name string, config interface{}) error {
	configDefine := s.configs[name]
	if configDefine == nil {
		return errors.New("Config Not Found ")
	}

	cfg := s.config.cfgs[name]
	if cfg == nil {
		config = configDefine.Default()
		return nil
	}

	body, _ := yaml.Marshal(cfg)
	return yaml.Unmarshal(body, config)
}

func (s *AppService) GetService(name string) IService {
	for k, v := range s.services {
		if v.ServiceName() == name {
			return s.services[k]
		}
	}

	return nil
}

type BaseService struct {
	IService
}

func (s *BaseService) Init(app ISptty) error {
	return nil
}

func (s *BaseService) Release() {}

func (s *BaseService) Enable() bool {
	return true
}

func (s *BaseService) ServiceName() string {
	return ""
}
