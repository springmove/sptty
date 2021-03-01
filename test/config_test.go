package test

import (
	"os"
	"path"
	"testing"

	"github.com/linshenqi/sptty"
)

// type cfg1 struct {
// 	Key1 string        `yaml:"key1"`
// 	Key2 int           `yaml:"key2"`
// 	Key3 time.Duration `yaml:"key3_wefli2"`
// }

// func (c1 *cfg1) ConfigName() string {
// 	return ""
// }

// func (c1 *cfg1) Validate() error {
// 	return nil
// }

// type cfg2 struct {
// 	Key4 []int             `yaml:"key4"`
// 	Key5 map[string]string `yaml:"key5"`
// }

// func (c2 *cfg2) ConfigName() string {
// 	return ""
// }

// func (c2 *cfg2) Validate() error {
// 	return nil
// }

func getSrv() *sptty.ConfigService {
	return &sptty.ConfigService{}
}

func TestConfig(t *testing.T) {
	// dir, _ := os.Getwd()
	// conf := path.Join(dir, "config.yml")

	// cfgs := sptty.Configs{
	// 	// &cfg1{},
	// 	//&cfg2{},
	// }

	// f, _ := os.Open(conf)
	// defer f.Close()

	// content, _ := ioutil.ReadAll(f)
	// _ = yaml.Unmarshal(content, &cfgs)

	// var cfg1 *cfg1
	// body, _ := yaml.Marshal(cfgs[0])
	// _ = yaml.Unmarshal(body, cfg1)
	// fmt.Println(cfg1)

	// cfg2 := cfg2{}
	// body, _ = yaml.Marshal(cfgs["cfg2"])
	// yaml.Unmarshal(body, &cfg2)
	// fmt.Println(cfg2)

	srv := getSrv()
	dir, _ := os.Getwd()
	srv.SetConfPath(path.Join(dir, "config.yml"))

	os.Setenv("sptty.cfg2.key4.0", "33")
	os.Setenv("sptty.cfg1.key3_wefli2", "10")

	_ = srv.Init(nil)
}
