package test

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func slog(tags ...string) {
	log.WithField("tag", tags).Info("wef")
}
func TestLog(t *testing.T) {
	//dir, _ := os.Getwd()
	//logOutput := path.Join(dir, "%Y%m%d%S.log")
	//logf, _ := rl.New(
	//	logOutput,
	//	rl.WithMaxAge(10 * time.Second),
	//	rl.WithRotationTime(5 * time.Second),
	//)

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	//log.WithField("tag", nil).Info("wef")
	//log.WithField("tag", "s1").Info("wefawef")
	//log.WithField("tag", "s2").Error("err234234")
	//log.WithField("tag", "s2").Debug("awef")

	slog("3r", "e")
}
