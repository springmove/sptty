package sptty

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	ModelServiceName = "model"
)

type ModelConfig struct {
	Enable  bool   `yaml:"enable"`
	Source  string `yaml:"source"`
	Name    string `yaml:"name"`
	User    string `yaml:"user"`
	Pwd     string `yaml:"pwd"`
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

func (c *ModelConfig) ConfigName() string {
	return ModelServiceName
}

func (c *ModelConfig) Validate() error {
	return nil
}

func (c *ModelConfig) Default() interface{} {
	return &ModelConfig{
		Enable: false,
	}
}

type ModelService struct {
	db *gorm.DB
}

func (s *ModelService) getConnStr(cfg *ModelConfig) string {
	connStr := ""

	switch cfg.Source {
	case "postgres":
		connStr = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s connect_timeout=%d sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Name,
			cfg.Pwd,
			cfg.Timeout)

	default:
		return connStr
	}

	return connStr
}

func (s *ModelService) Init(app ISptty) error {

	s.db = nil

	cfg := ModelConfig{}
	err := app.GetConfig(s.ServiceName(), &cfg)
	if err != nil {
		return err
	}

	if !cfg.Enable {
		Log(InfoLevel, fmt.Sprintf("%s Service Is Disabled", s.ServiceName()), s.ServiceName())
		return nil
	}

	_db, err := gorm.Open(cfg.Source, s.getConnStr(&cfg))
	if err != nil {
		return err
	}

	s.db = _db

	return nil
}

func (s *ModelService) AddModel(values interface{}) {
	if s.db != nil {
		s.db.AutoMigrate(values)
	}
}

func (s *ModelService) DB() *gorm.DB {
	return s.db
}

func (s *ModelService) Release() {
	if s.db != nil {
		_ = s.db.Close()
	}
}

func (s *ModelService) Enable() bool {
	return true
}

func (s *ModelService) ServiceName() string {
	return ModelServiceName
}

type SimpleModelBase struct {
	ID      string     `gorm:"size:32;primary_key" json:"id"`
	Created *time.Time `json:"created,omitempty"`
	Deleted bool       `json:"-"`
}
