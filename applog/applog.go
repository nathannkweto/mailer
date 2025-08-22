package applog

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Fields = logrus.Fields

var Log = logrus.New()

// InitLogger configures a JSON logger, level from env (debug/info/warn/error)
func InitLogger(level string) {
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{})

	switch strings.ToLower(level) {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "warn", "warning":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
}
