package sptty

import (
	fr "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	LogServiceName = "LOG"
	LogStdout      = "STDOUT"
	LogTag         = "tag"

	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARNING"
	ErrorLevel = "ERROR"
	FatalLevel = "FATAL"
)

type LogLevel string

var LogLevels = map[string]log.Level{
	DebugLevel: log.DebugLevel,
	InfoLevel:  log.InfoLevel,
	WarnLevel:  log.WarnLevel,
	ErrorLevel: log.ErrorLevel,
	FatalLevel: log.FatalLevel,
}

type LogConfig struct {
	File   string        `yaml:"file"`
	Level  string        `yaml:"level"`
	MaxAge time.Duration `yaml:"max_age"`
	Rotate time.Duration `yaml:"rotate"`
}

type LogService struct {
	cfg *LogConfig
}

func (s *LogService) Init(app Sptty) error {

	err := app.GetConfig(s.ServiceName(), s.cfg)
	if err != nil {
		s.cfg = s.defaultConfig()
	}

	s.setupLog()
	return nil
}

func (s *LogService) setupLog() {
	if s.cfg.File != LogStdout {
		logf, _ := fr.New(
			s.cfg.File,
			fr.WithMaxAge(s.cfg.MaxAge),
			fr.WithRotationTime(s.cfg.Rotate),
		)
		log.SetOutput(logf)
	} else {
		log.SetOutput(os.Stdout)
	}

	level, exist := LogLevels[s.cfg.Level]
	if !exist {
		level = log.DebugLevel
	}

	log.SetLevel(level)
}

func (s *LogService) Log(level LogLevel, msg string, tags ...string) {
	switch level {
	case DebugLevel:
		log.WithField(LogTag, tags).Debug(msg)
	case InfoLevel:
		log.WithField(LogTag, tags).Info(msg)
	case WarnLevel:
		log.WithField(LogTag, tags).Warn(msg)
	case ErrorLevel:
		log.WithField(LogTag, tags).Error(msg)
	case FatalLevel:
		log.WithField(LogTag, tags).Fatal(msg)
	default:
		log.WithField(LogTag, tags).Debug(msg)
	}
}

func (s *LogService) defaultConfig() *LogConfig {
	return &LogConfig{
		File:   "STDOUT",
		Level:  "DEBUG",
		MaxAge: 2160 * time.Hour,
		Rotate: 24 * time.Hour,
	}
}

func (s *LogService) Release() {

}

func (s *LogService) Enable() bool {
	return true
}

func (s *LogService) ServiceName() string {
	return LogServiceName
}
