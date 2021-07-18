package sptty

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ModelServiceName = "model"

	DBPostgres = "postgres"
	DBMysql    = "mysql"

	DefaultMysqlCharset = "utf8mb4"
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

	// for mysql
	Charset string `yaml:"charset"`
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
	case DBPostgres:
		// ex: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
		connStr = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s connect_timeout=%d sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Name,
			cfg.Pwd,
			cfg.Timeout)

	case DBMysql:
		if cfg.Charset == "" {
			cfg.Charset = DefaultMysqlCharset
		}

		// ex: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User,
			cfg.Pwd,
			cfg.Host,
			cfg.Port,
			cfg.Name,
			cfg.Charset)

	default:
		return connStr
	}

	return connStr
}

func (s *ModelService) Init(app ISptty) error {

	var err error
	s.db = nil

	cfg := ModelConfig{}
	if err = app.GetConfig(s.ServiceName(), &cfg); err != nil {
		return err
	}

	if !cfg.Enable {
		Log(InfoLevel, fmt.Sprintf("%s Service Is Disabled", s.ServiceName()), s.ServiceName())
		return nil
	}

	switch cfg.Source {
	case DBPostgres:
		s.db, err = gorm.Open(postgres.Open(s.getConnStr(&cfg)), &gorm.Config{})
		if err != nil {
			return err
		}

	case DBMysql:
		s.db, err = gorm.Open(mysql.Open(s.getConnStr(&cfg)), &gorm.Config{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Source Not Supported")
	}

	return nil
}

func (s *ModelService) AddModel(values interface{}) {
	if s.db != nil {
		if err := s.db.AutoMigrate(values); err != nil {
			Log(ErrorLevel, fmt.Sprintf("AutoMigrate Failed: %s", err.Error()), s.ServiceName())
		}
	}
}

func (s *ModelService) DB() *gorm.DB {
	return s.db
}

func (s *ModelService) Release() {
	if s.db != nil {
		_ = s.db
	}
}

func (s *ModelService) Enable() bool {
	return true
}

func (s *ModelService) ServiceName() string {
	return ModelServiceName
}

type SimpleModelBase struct {
	ID        string    `gorm:"size:32;primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Deleted   bool      `json:"deleted,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (s *SimpleModelBase) Init() *SimpleModelBase {
	s.ID = GenerateUID()
	s.CreatedAt = time.Now().UTC()
	s.Deleted = false

	return s
}

func (s *SimpleModelBase) Serialize() *SimpleModelBase {
	s.CreatedAt = s.CreatedAt.UTC()
	s.UpdatedAt = s.UpdatedAt.UTC()

	return s
}

func UpdateModel(db *gorm.DB, model interface{}) error {
	if err := db.Save(model).Error; err != nil {
		return err
	}

	return nil
}
