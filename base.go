package sptty

type Configs []IConfig
type Services []IService

type SerivcesHandler func(ISptty) Services
type ConfigsHandler func(ISptty) Configs

type ISptty interface {
	GetService(name string) IService
	GetConfig(name string, config IConfig) error
	AddServices(handler SerivcesHandler)
	AddConfigs(handler ConfigsHandler)
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
	Default() IConfig
}

type IServices interface {
	Services() Services
	Configs() Configs
}
