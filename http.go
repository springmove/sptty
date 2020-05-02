package sptty

import (
	"io/ioutil"
	"time"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gopkg.in/resty.v1"
)

const (
	BaseApiRoute    = "/api"
	HttpServiceName = "http"
)

type HttpConfig struct {
	Addr   string `yaml:"addr"`
	Tag    string `yaml:"tag"`
	ApiDoc string `yaml:"api_doc"`
}

func (c *HttpConfig) ConfigName() string {
	return HttpServiceName
}

func (c *HttpConfig) Validate() error {
	return nil
}

func (c *HttpConfig) Default() interface{} {
	return &HttpConfig{
		Addr:   "8080",
		ApiDoc: "",
	}
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
	Timeout      time.Duration     `yaml:"timeout"`
	Headers      map[string]string `yaml:"headers"`
	PushInterval time.Duration     `yaml:"push_interval"`
	MaxRetry     int               `yaml:"max_retry"`
}

type HttpService struct {
	app   *iris.Application
	party iris.Party
	cfg   HttpConfig
}

func DefaultHttpClientConfig() *HttpClientConfig {
	return &HttpClientConfig{
		Timeout:      8 * time.Second,
		PushInterval: 1 * time.Second,
		MaxRetry:     3,
		Headers:      map[string]string{},
	}
}

func CreateHttpClient(cfg *HttpClientConfig) *resty.Client {
	client := resty.New()

	client.SetRESTMode()
	client.SetTimeout(cfg.Timeout)
	client.SetContentLength(true)
	client.SetHeaders(cfg.Headers)
	client.
		SetRetryCount(cfg.MaxRetry).
		SetRetryWaitTime(cfg.PushInterval).
		SetRetryMaxWaitTime(20 * time.Second)

	return client
}

func (s *HttpService) SetOptions() {
	tag := s.cfg.Tag
	if tag == "" {
		tag = BaseApiRoute
	}

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
	})

	s.party = s.app.Party(tag, crs).AllowMethods(iris.MethodOptions)
}

func (s *HttpService) Init(app Sptty) error {
	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return err
	}

	s.AddRoute("GET", "healthz", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
	})

	s.AddRoute("GET", "/healthz", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
	})

	s.AddRoute("GET", "/apidoc", func(ctx iris.Context) {
		ctx.Header("content-type", "application/json")
		f, err := ioutil.ReadFile(s.cfg.ApiDoc)
		if err != nil {
			ctx.StatusCode(iris.StatusNoContent)
			return
		}
		_, _ = ctx.Write(f)
	})

	return s.app.Run(iris.Addr(s.cfg.Addr), iris.WithoutInterruptHandler)
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
