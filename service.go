package sptty

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v2"
)

var appService *AppService = nil

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

	Log(InfoLevel, "Now Sptting~", "sptty")

	return nil
}

func (s *AppService) release() {
	for _, v := range s.services {
		v.Release()
	}
}

func (s *AppService) Sptting() {
	defer func() {
		s.release()
	}()

	if err := s.init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func (s *AppService) validateConfigs() error {
	for _, v := range s.configs {
		if err := v.Validate(); err != nil {
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

func (s *AppService) GetConfig(name string, config IConfig) error {
	configDefine := s.configs[name]
	if configDefine == nil {
		return fmt.Errorf("Config Not Found ")
	}

	cfg := s.config.cfgs[name]
	if cfg == nil {
		config = configDefine.Default()
		return nil
	}

	body, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

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

func (s *AppService) AddServices(handler SerivcesHandler) {
	s.services = handler(s)
}

func (s *AppService) AddConfigs(handler ConfigsHandler) {
	configs := handler(s)
	for k, v := range configs {
		s.configs[v.ConfigName()] = configs[k]
	}
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
