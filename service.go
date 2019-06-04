package sptty

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/mitchellh/mapstructure"
)

var appService *AppService = nil

func GetApp() *AppService {
	if appService == nil {
		appService = &AppService{
			http: &HttpService{
				app: iris.New(),
			},
			model:    &ModelService{},
			config:   &Config{},
			services: map[string]Service{},
		}

		appService.http.(*HttpService).SetOptions()
	}

	return appService
}

type AppService struct {
	services map[string]Service
	http     Service
	model    Service
	config   Service
	Sptty
}

func (bs *AppService) init() {
	if bs.config.Init(bs) != nil {
		return
	}

	if bs.model.Init(bs) != nil {
		return
	}

	for _, v := range bs.services {
		if v.Init(bs) != nil {
			return
		}
	}

	bs.http.Init(bs)
}

func (bs *AppService) release() {
	for _, v := range bs.services {
		v.Release()
	}

	bs.http.Release()
}

func (bs *AppService) Sptting() {
	bs.init()
	bs.release()
}

func (bs *AppService) AddServices(services map[string]Service) {
	bs.services = services
}

func (bs *AppService) AddConfigs(cfgs SpttyConfig) {
	bs.config.(*Config).AddConfigs(cfgs)
}

func (bs *AppService) ConfFromFile(conf string) {
	bs.config.(*Config).SetConfPath(conf)
}

func (bs *AppService) cfg() SpttyConfig {
	config := bs.config.(*Config)
	return config.cfg
}

func (bs *AppService) GetConfig(name string, config interface{}) error {
	cfg := bs.cfg()[name]
	if cfg == nil {
		return nil
	}

	return mapstructure.Decode(cfg, config)
}

func (bs *AppService) AddRoute(method string, route string, handler context.Handler) {
	bs.http.(*HttpService).AddRoute(method, route, handler)
}

func (bs *AppService) AddModel(values interface{}) {
	bs.model.(*ModelService).AddModel(values)
}

func (bs *AppService) Http() Service {
	return bs.http
}

func (bs *AppService) Model() Service {
	return bs.model
}

func (bs *AppService) GetService(name string) Service {
	return bs.services[name]
}

func (bs *AppService) RegistService(name string, service Service) {
	bs.services[name] = service
}
