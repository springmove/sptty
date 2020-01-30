package sptty

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/resty.v1"
	"time"
)

const (
	BaseApiRoute    = "/api/v1"
	HttpServiceName = "http"
)

type HttpConfig struct {
	Addr string
}

type Route struct {
	RouteType   string
	Method      string
	Pattern     string
	HandlerFunc context.Handler
}

type CorsConfig struct {
	AllowedOrigins   []string `yaml:"allowed-origins"`
	AllowCredentials bool     `yaml:"allow-credentials"`
	AllowedMethods   []string `yaml:"allowed-methods"`
}

type HttpClientConfig struct {
	Timeout      int               `yaml:"timeout"`
	Headers      map[string]string `yaml:"headers"`
	PushInterval int               `yaml:"push_interval"`
	MaxRetry     int               `yaml:"max_retry"`
}

type HttpService struct {
	app   *iris.Application
	party iris.Party
}

func CreateHttpClient(cfg *HttpClientConfig) *resty.Client {
	client := resty.New()

	client.SetRESTMode()
	client.SetTimeout(time.Duration(cfg.Timeout) * time.Second)
	client.SetContentLength(true)
	client.SetHeaders(cfg.Headers)
	client.
		SetRetryCount(cfg.MaxRetry).
		SetRetryWaitTime(time.Duration(cfg.PushInterval) * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	return client
}

func (s *HttpService) SetOptions() {
	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Next()
	}

	s.party = s.app.Party(BaseApiRoute, crs).AllowMethods(iris.MethodOptions)
}

func (s *HttpService) Init(app Sptty) error {
	cfg := HttpConfig{}
	app.GetConfig("s", &cfg)

	s.app.Run(iris.Addr(cfg.Addr), iris.WithoutInterruptHandler)

	return nil
}

func (s *HttpService) Release() {

}

func (s *HttpService) Enable() bool {
	return true
}

func (s *HttpService) AddRoute(method string, route string, handler context.Handler) {
	s.party.Handle(method, route, handler)
}

func (s *HttpService) ServiceName() string {
	return HttpServiceName
}
