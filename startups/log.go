package startups

import (
	"os"

	"github.com/sirupsen/logrus"
)

func StartupLog() *logrus.Logger {

	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(dir+"/hots.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	log.SetOutput(file)
	log.SetLevel(logrus.InfoLevel)
	log.SetReportCaller(true)

	return log
}
