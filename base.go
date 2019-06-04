package sptty

import (
	"github.com/kataras/iris/context"
)

type SpttyConfig map[string]interface{}

type Sptty interface {
	GetService(name string) Service
	AddServices(services map[string]Service)
	AddConfigs(cfg SpttyConfig)
	GetConfig(name string, config interface{}) error
	Http() Service
	Model() Service
	AddRoute(method string, route string, handler context.Handler)
	AddModel(values interface{})
}

type Service interface {
	Init(app Sptty) error
	Release()
	Enable() bool
}
