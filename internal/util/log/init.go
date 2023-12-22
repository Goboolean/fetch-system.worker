package log

import (
	"os"

	"github.com/sirupsen/logrus"
)



func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		DisableTimestamp: false,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
}