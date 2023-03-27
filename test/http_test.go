package test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/springmove/sptty"
)

type testHttp struct {
}

type pkg struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
}

func (s *testHttp) Init(app sptty.ISptty) error {
	app.AddRoute("POST", "/auth", func(ctx iris.Context) {
		var req pkg
		if err := ctx.ReadJSON(&req); err != nil {
			fmt.Println(err.Error())
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		fmt.Println(req)
	})

	app.AddRoute("POST", "/users", func(ctx iris.Context) {
		var req interface{}
		if err := ctx.ReadJSON(&req); err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(req)
	})

	app.AddRoute("GET", "/user", func(ctx iris.Context) {

		_ = sptty.SimpleResponse(ctx, iris.StatusOK, map[string]interface{}{
			"awef": 23,
		})
	})

	return nil
}

func (s *testHttp) Release() {
}

func (s *testHttp) Enable() bool {
	return true
}

func (s *testHttp) ServiceName() string {
	return "testHttp"
}

func TestHttp(t *testing.T) {
	dir, _ := os.Getwd()
	conf := path.Join(dir, "config.yml")

	app := sptty.GetApp()
	app.ConfFromFile(conf)

	app.AddServices(sptty.Services{
		&testHttp{},
	})

	app.Sptting()
}
