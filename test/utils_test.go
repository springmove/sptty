package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestUID(t *testing.T) {
	// uids := []string{}
	// for i := 0; i < 7; i++ {
	// 	uids = append(uids, sptty.GenerateUID())
	// }

	td, _ := time.ParseDuration("24h")
	var ti interface{} = td

	fmt.Println(reflect.TypeOf(ti).String() == "time.Duration")
}
