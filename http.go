package sptty

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/resty.v1"
	"time"
)

const (
	BASE_API_ROUTE = "/api/v1"
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
	Service
	//cors CorsConfig
	//routes []Route
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

func (http *HttpService) SetOptions() {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	http.party = http.app.Party(BASE_API_ROUTE, crs).AllowMethods(iris.MethodOptions)
}

func (http *HttpService) Init(app Sptty) error {
	cfg := HttpConfig{}
	app.GetConfig("http", &cfg)

	http.app.Run(iris.Addr(cfg.Addr), iris.WithoutInterruptHandler)

	return nil
}

func (http *HttpService) Release() {

}

func (http *HttpService) Enable() bool {
	return true
}

func (http *HttpService) AddRoute(method string, route string, handler context.Handler) {
	http.party.Handle(method, route, handler)
}
