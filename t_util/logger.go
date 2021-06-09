package t_util

import (
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {
	logLevel := flag.String("log-level", "debug", "Log level")
	logFile := flag.String("log-file", "", "Log file")
	flag.Parse()

	*logLevel = strings.ToLower(*logLevel)

	switch *logLevel {
	default:
		Log.Fatal("Bad input for 'log-level' flag")
	case "debug":
		Log.Level = logrus.DebugLevel
	case "info":
		Log.Level = logrus.InfoLevel
	case "warn", "warning":
		Log.Level = logrus.WarnLevel
	case "error":
		Log.Level = logrus.ErrorLevel
	case "panic":
		Log.Level = logrus.PanicLevel
	case "fatal":
		Log.Level = logrus.FatalLevel
	}
	// logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			Log.Out = file
		} else {
			Log.WithFields(logrus.Fields{
				"file": *logFile,
			}).Info("Failed to log to file, using default stderr")
		}
	} else {
		Log.Info("No logfile provided, using default stderr")
	}
	Log.Infof("Logging on %v level", Log.Level)
}
