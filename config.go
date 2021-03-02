package sptty

import (
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	ConfigServiceName = "config"
	ConfigEnvPrefix   = "sptty"
	ConfigEnvKeyDiv   = "_"
)

type ConfigService struct {
	confPath string
	cfgs     map[interface{}]interface{}
}

func (s *ConfigService) Init(app ISptty) error {

	f, err := os.Open(s.confPath)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &s.cfgs)
	if err != nil {
		return err
	}

	s.patchConfigsWithEnvs()
	return nil
}

func getTargetEnvs() map[string]string {
	rt := map[string]string{}
	for _, env := range os.Environ() {
		vals := strings.SplitN(env, "=", 2)
		if strings.Contains(vals[0], ConfigEnvPrefix) {
			rt[strings.TrimSpace(vals[0])] = strings.TrimSpace(vals[1])
		}
	}

	return rt
}

func (s *ConfigService) patchConfigWithEnv(key string, value string) {
	keys := strings.Split(key, ConfigEnvKeyDiv)
	l := len(keys)
	step := "obj"
	var obj map[interface{}]interface{} = s.cfgs
	var arr []interface{}
	for i := 1; i < l; i++ {

		k := keys[i]
		switch step {
		case "obj":
			val, exist := obj[k]
			if !exist {
				return
			}
			typ := reflect.TypeOf(val).Kind()
			switch typ {
			case reflect.Map:
				obj = val.(map[interface{}]interface{})
				step = "obj"
			case reflect.Slice:
				arr = val.([]interface{})
				step = "arr"
			default:
				// value
				switch typ {
				case reflect.Bool:
					obj[k], _ = strconv.ParseBool(value)
				case reflect.Int:
					obj[k], _ = strconv.ParseInt(value, 10, 32)
				case reflect.Float32:
					obj[k], _ = strconv.ParseFloat(value, 32)
				case reflect.String:
					obj[k] = value
				}
			}
		case "arr":
			index, err := strconv.Atoi(k)
			if err != nil {
				return
			}

			if index >= len(arr) {
				return
			}

			val := arr[index]
			typ := reflect.TypeOf(val).Kind()
			switch typ {
			case reflect.Map:
				obj = val.(map[interface{}]interface{})
				step = "obj"
			case reflect.Slice:
				arr = val.([]interface{})
				step = "arr"
			default:
				// value
				switch typ {
				case reflect.Bool:
					arr[index], _ = strconv.ParseBool(value)
				case reflect.Int:
					arr[index], _ = strconv.ParseInt(value, 10, 32)
				case reflect.Float32:
					arr[index], _ = strconv.ParseFloat(value, 32)
				case reflect.String:
					arr[index] = value
				}
			}
		}
	}
}

func (s *ConfigService) patchConfigsWithEnvs() {
	envs := getTargetEnvs()
	for k, v := range envs {
		s.patchConfigWithEnv(k, v)
	}
}

func (s *ConfigService) Release() {

}

func (s *ConfigService) Enable() bool {
	return true
}

func (s *ConfigService) SetConfPath(conf string) {
	s.confPath = conf
}

func (s *ConfigService) ServiceName() string {
	return ConfigServiceName
}

type BaseConfig struct {
	IConfig
}

func (s *BaseConfig) ConfigName() string {
	return ""
}

func (s *BaseConfig) Validate() error {
	return nil
}

func (s *BaseConfig) Default() interface{} {
	return &BaseConfig{}
}
