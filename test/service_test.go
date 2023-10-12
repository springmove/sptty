package test

import (
	"fmt"
	"testing"

	"github.com/springmove/sptty"
)

type ConfigApp struct {
	sptty.BaseConfig

	KK int
}

type SpttyApp struct {
	MyService IMyService

	Config ConfigApp
}

type MyService struct {
	sptty.BaseService
}

func (s *MyService) ServiceName() string {
	return "MyService"
}

func (s *MyService) Init(app sptty.ISptty) error {
	fmt.Println("MyService.Init")
	return nil
}

func (s *MyService) F1() {
	fmt.Println("MyService.F1")
}

type IMyService interface {
	F1()
}

func TestService(t *testing.T) {
	// app := SpttyApp{
	// 	MyService: &MyService{},

	// 	Config: ConfigApp{},
	// }
}
