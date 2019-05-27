package sptty

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func (s *ModelService) Init(app Sptty) error {

	s.db = nil

	cfg := ModelConfig{}
	err := app.GetConfig("model", &cfg)
	if err != nil {
		return err
	}

	if !cfg.Enable {
		return nil
	}

	_db, err := gorm.Open(cfg.Source, s.getConnStr(&cfg))
	if err != nil {
		return err
	}

	s.db = _db

	return nil
}

func (s *ModelService) AddModel(m interface{}) {
	if s.db != nil {
		s.db.AutoMigrate(m)
	}
}

func (s *ModelService) DB() *gorm.DB {
	return s.db
}

func (s *ModelService) Release() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s ModelService) Enable() bool {
	return true
}
