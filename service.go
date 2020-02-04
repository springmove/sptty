package sptty

import (
	"errors"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/yaml.v2"
)

var appService *AppService = nil

func GetApp() *AppService {
	if appService == nil {
		appService = &AppService{
			http: &HttpService{
				app: iris.New(),
			},
			model:    &ModelService{},
			config:   &ConfigService{},
			log:      &LogService{},
			services: map[string]Service{},
			configs:  map[string]Config{},
		}

		appService.http.SetOptions()
	}

	return appService
}

func Log(level LogLevel, msg string, tags ...string) {
	app := GetApp()
	log := app.GetService(LogServiceName).(*LogService)
	log.Log(level, msg, tags...)
}

type AppService struct {
	services map[string]Service
	configs  map[string]Config
	http     *HttpService
	model    *ModelService
	config   *ConfigService
	log      *LogService
	Sptty
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

	Log(InfoLevel, fmt.Sprintf("init service: %s", s.model.ServiceName()), s.model.ServiceName())
	if err := s.model.Init(s); err != nil {
		return err
	}

	for _, v := range s.services {
		Log(InfoLevel, fmt.Sprintf("init service: %s", v.ServiceName()), v.ServiceName())
		if err := v.Init(s); err != nil {
			return err
		}
	}

	if err := s.http.Init(s); err != nil {
		return err
	}

	return nil
}

func (s *AppService) release() {
	for _, v := range s.services {
		v.Release()
	}

	s.http.Release()
}

func (s *AppService) Sptting() {
	if s.init() != nil {
		return
	}

	s.release()
}

func (s *AppService) AddServices(services Services) {
	for k, v := range services {
		s.services[v.ServiceName()] = services[k]
	}
}

func (s *AppService) AddConfigs(cfgs Configs) {
	for k, v := range cfgs {
		s.configs[v.ConfigName()] = cfgs[k]
	}
}

func (s *AppService) validateConfigs() error {
	for _, v := range s.configs {
		err := v.Validate()
		if err != nil {
			fmt.Printf("Config Error: %s", err.Error())
			return err
		}
	}

	return nil
}

func (s *AppService) ConfFromFile(conf string) {
	s.config.SetConfPath(conf)
}

func (s *AppService) GetConfig(name string, config interface{}) error {
	cfg := s.config.cfgs[name]
	if cfg == nil {
		return errors.New(" Config Not Found")
	}

	body, _ := yaml.Marshal(cfg)
	return yaml.Unmarshal(body, config)
}

func (s *AppService) AddRoute(method string, route string, handler context.Handler) {
	s.http.AddRoute(method, route, handler)
}

func (s *AppService) AddModel(values interface{}) {
	s.model.AddModel(values)
}

func (s *AppService) Http() Service {
	return s.http
}

func (s *AppService) Model() Service {
	return s.model
}

func (s *AppService) GetService(name string) Service {
	return s.services[name]
}

func (s *AppService) RegistService(name string, service Service) {
	s.services[name] = service
}
