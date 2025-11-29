package test

import (
	"fmt"
	"testing"

	"github.com/springmove/sptty"
)

func TestUID(t *testing.T) {
	// uids := []string{}
	// for i := 0; i < 7; i++ {
	// 	uids = append(uids, sptty.GenerateUID())
	// }

	// td, _ := time.ParseDuration("24h")
	// var ti interface{} = td

	// fmt.Println(reflect.TypeOf(ti).String() == "time.Duration")

	// _, mime, err := sptty.GetUrlImage("https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTI7YP85icEWiaYmNaYFUrvNdjiaCcFMxO0Wj9dH7QfX5ZtibHPW50HDsibLoHApUic5EALnja3vnn0uKziaQ/132")
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// 	return
	// }

	fmt.Println(sptty.Sha1("93102ba73b19e781468b157909274353a8ceb6fa"))
}
