package sptty

import (
	"errors"
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
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
			http: &HttpService{
				app: iris.New(),
			},
			model:    &ModelService{},
			config:   &ConfigService{},
			log:      &LogService{},
			i18n:     &I18NService{},
			services: Services{},
			configs: map[string]Config{
				HttpServiceName:  &HttpConfig{},
				ModelServiceName: &ModelConfig{},
				LogServiceName:   &LogConfig{},
				I18NServiceName:  &I18NConfig{},
			},
		}

		appService.http.SetOptions()
	}

	return appService
}

func Log(level LogLevel, msg string, tags ...string) {
	app := GetApp()
	app.log.Log(level, msg, tags...)
}

func I18NValue(name string, lang string) string {
	app := GetApp()
	return app.i18n.get(name, lang)
}

type AppService struct {
	services Services
	configs  map[string]Config
	http     *HttpService
	model    *ModelService
	config   *ConfigService
	log      *LogService
	i18n     *I18NService
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

	if err := s.i18n.Init(s); err != nil {
		return err
	}

	Log(InfoLevel, fmt.Sprintf("Init Service: %s", s.model.ServiceName()), s.model.ServiceName())
	if err := s.model.Init(s); err != nil {
		Log(ErrorLevel, fmt.Sprintf("Init Service %s failed: %s", s.model.ServiceName(), err.Error()), s.model.ServiceName())
		return err
	}

	for _, v := range s.services {
		Log(InfoLevel, fmt.Sprintf("Init Service: %s", v.ServiceName()), v.ServiceName())
		if err := v.Init(s); err != nil {
			Log(ErrorLevel, fmt.Sprintf("Init Service %s failed: %s", v.ServiceName(), err.Error()), v.ServiceName())
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
	if err := s.init(); err != nil {
		fmt.Println(err.Error())
		return
	}

	s.release()
}

func (s *AppService) AddServices(services Services) {
	s.services = services
}

func (s *AppService) AddConfigs(configs Configs) {
	for k, v := range configs {
		s.configs[v.ConfigName()] = configs[k]
	}
}

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
	for k, v := range s.services {
		if v.ServiceName() == name {
			return s.services[k]
		}
	}

	return nil
}

type BaseService struct {
	Service
}

func (s *BaseService) Init(app Sptty) error {
	return nil
}

func (s *BaseService) Release() {}

func (s *BaseService) Enable() bool {
	return true
}

func (s *BaseService) ServiceName() string {
	return ""
}
