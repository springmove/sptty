package sptty

type Configs []IConfig
type Services []IService

type ISptty interface {
	// GetService(name string) IService
	// AddServices(services Services)
	// AddConfigs(configs Configs)
	// GetConfig(name string, config interface{}) error
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
