package sptty

import (
	"github.com/kataras/iris/context"
)

type SpttyConfig map[string]interface{}

type Sptty interface {
	//Config() SpttyConfig
	GetService(name string) Service
	AddServices(services map[string]Service)
	AddConfigs(cfg SpttyConfig)
	GetConfig(name string, config interface{})
	Http() Service
	AddRoute(method string, route string, handler context.Handler)
}

type Service interface {
	Init(service Sptty) error
	Release()
}
