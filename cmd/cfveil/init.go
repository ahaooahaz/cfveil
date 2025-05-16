package main

import (
	"github.com/ahaooahaz/cfveil/cmd/cfveil/python"

	"github.com/sirupsen/logrus"
)

func init() {
	err := initEnv()
	if err != nil {
		panic(err.Error())
	}

	rootCmd.AddCommand(python.Cmd)

}

func initEnv() (err error) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return
}
