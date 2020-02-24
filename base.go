package sptty

import (
	"github.com/kataras/iris/v12/context"
)

type Configs []Config
type Services []Service

type Sptty interface {
	GetService(name string) Service
	AddServices(services Services)
	AddConfigs(configs Configs)
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
	ServiceName() string
}

type Config interface {
	ConfigName() string
	Validate() error
	Default() interface{}
}
