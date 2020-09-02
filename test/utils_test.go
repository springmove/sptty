package test

import (
	"fmt"
	"testing"

	"github.com/linshenqi/sptty"
)

func TestUID(t *testing.T) {
	uids := []string{}
	for i := 0; i <7; i++ {
		uids = append(uids, sptty.GenerateUID())
	}
	fmt.Println(uids)
}
