package sptty

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	I18NServiceName = "i18n"

	ZH = "zh" // 中文
	EN = "en" // 英语
	FR = "fr" // 法语
	DE = "de" // 德语
	JA = "ja" // 日语
	KO = "ko" // 韩语
	ES = "es" // 西班牙语
	PT = "pt" // 葡萄牙语
	RU = "ru" // 俄语
	IT = "it" // 意大利语
	AR = "ar" // 阿拉伯语
	MS = "ms" // 马来语
	TH = "th" // 泰语
	AF = "af" // 南非语
)

type I18NConfig struct {
	File string `yaml:"file"`
}

func (c *I18NConfig) ConfigName() string {
	return I18NServiceName
}

func (c *I18NConfig) Validate() error {
	return nil
}

func (c *I18NConfig) Default() interface{} {
	return &I18NConfig{
		File: "",
	}
}

type I18NService struct {
	cfg   I18NConfig
	trans map[string]map[string]string
}

func (s *I18NService) Init(app Sptty) error {
	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return err
	}

	return s.load()
}

func (s *I18NService) Release() {
}

func (s *I18NService) Enable() bool {
	return true
}

func (s *I18NService) ServiceName() string {
	return I18NServiceName
}

func (s *I18NService) load() error {
	if s.cfg.File == "" {
		return nil
	}

	f, err := os.Open(s.cfg.File)
	defer f.Close()

	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &s.trans)
}

func (s *I18NService) get(name string, lang string) string {
	target, exist := s.trans[name]
	if !exist {
		return name
	}

	langValue, exist := target[lang]
	if !exist {
		return name
	}

	return langValue
}
