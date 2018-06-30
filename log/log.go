package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/moxitrel/golib/cfg"
)

func init() {
	logrus.SetLevel(cfg.DebugLevel)
	//logrus.SetOutput(os.Stdout)
}

func withFunName(f func(string, ...interface{})) func(string, ...interface{}) {
	return func(format string, args ...interface{}) {
		format = "%s: " + format
		args = append([]interface{}{CallerName(1)}, args...)
		f(format, args...)
	}
}

var Debug = withFunName(logrus.Debugf)
var Info = withFunName(logrus.Infof)
var Warn = withFunName(logrus.Warnf)
var Error = withFunName(logrus.Errorf)
var Fatal = withFunName(logrus.Fatalf)
