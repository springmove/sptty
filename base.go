package sptty

import (
	"github.com/kataras/iris/v12/context"
)

type Configs []IConfig
type Services []IService

type ISptty interface {
	GetService(name string) IService
	AddServices(services Services)
	AddConfigs(configs Configs)
	GetConfig(name string, config interface{}) error
	Http() IService
	Model() IService
	AddRoute(method string, route string, handler context.Handler)
	AddModel(values interface{})
}

type IService interface {
	Init(app ISptty) error
	Release()
	Enable() bool
	ServiceName() string
}

type IConfig interface {
	ConfigName() string
	Validate() error
	Default() interface{}
}
