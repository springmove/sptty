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
			config:   &Config{},
			services: map[string]Service{},
		}

		appService.Http().(*HttpService).SetOptions()
	}

	return appService
}

type AppService struct {
	services map[string]Service
	http     Service
	config   Service
	Sptty
}

func (bs *AppService) init() {
	if bs.config.Init(bs) != nil {
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

func (bs *AppService) SetConf(conf string) {
	bs.config.(*Config).SetConf(conf)
}

func (bs *AppService) cfg() SpttyConfig {
	config := bs.config.(*Config)
	return config.cfg
}

func (bs *AppService) GetConfig(name string, config interface{}) {
	cfg := bs.cfg()[name]
	if cfg == nil {
		return
	}

	mapstructure.Decode(cfg, config)
}

func (bs *AppService) AddRoute(method string, route string, handler context.Handler) {
	http := bs.http.(*HttpService)
	http.AddRoute(method, route, handler)
}

func (bs *AppService) Http() Service {
	return bs.http
}

func (bs *AppService) GetService(name string) Service {
	return bs.services[name]
}

func (bs *AppService) RegistService(name string, service Service) {
	bs.services[name] = service
}
