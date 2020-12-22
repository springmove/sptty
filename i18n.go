package sptty

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
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
	Path string `yaml:"path"`
}

func (c *I18NConfig) ConfigName() string {
	return I18NServiceName
}

func (c *I18NConfig) Validate() error {
	return nil
}

func (c *I18NConfig) Default() interface{} {
	return &I18NConfig{
		Path: "",
	}
}

type I18NService struct {
	cfg   I18NConfig
	trans map[string]map[string]string
}

func (s *I18NService) Init(app ISptty) error {

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
	if s.cfg.Path == "" {
		return nil
	}

	files, err := ioutil.ReadDir(s.cfg.Path)
	if err != nil {
		return err
	}

	for _, v := range files {
		if err := s.loadFileContent(path.Join(s.cfg.Path, v.Name())); err != nil {
			continue
		}
	}

	return nil
}

func (s *I18NService) loadFileContent(filepath string) error {
	f, err := os.Open(filepath)
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

	return json.Unmarshal(content, &s.trans)
}

func (s *I18NService) get(name string, lang string) string {
	if s.trans == nil {
		return name
	}

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
