package main

import (
	"flag"
	"encoding/json"
	"os"
	"io/ioutil"
	"log"

	"github.com/tdecker91/gws"
	"github.com/sirupsen/logrus"
)

var port = flag.Int("port", 8030, "port to listen on")
var route = flag.String("route", "/", "route to listen for socket connections")
var level = flag.String("level", "info", "logging level. [debug, info, warn, error]")
var formatter = flag.String("formatter", "text", "format to output log messages as. [text, json]")
var timeFormat = flag.String("timeFormat", "2006-01-02 15:04:05", "text formatter only. format string to output time logs. Format as defined here https://golang.org/pkg/time/#Time.Format")
var format = flag.String("format", "[%L] %D: %M", "text formatter only. format for the output message. %L = log level. %D = Date. %M = Message")

func init() {
	flag.Parse()

	// Disable external packages from logging using log.Logger
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	switch *level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	switch *formatter {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		fmtr := &GowlFormatter{*format}
		logrus.SetFormatter(fmtr)
	}

	logrus.SetOutput(os.Stdout)
}

func main() {
	messages := make(chan gws.Message)

	server := gws.NewSocketServer(*port, *route)
	go server.Start(messages)

	logrus.Infof("GOWL listening for connections on ws://localhost:%d%s", *port, *route);


	for {
		m := <-messages

		switch t := m.Type; t {
		case gws.ClientMessage:
			var wm WebsocketMessage
			if err := json.Unmarshal(m.Data, &wm); err != nil {
				logrus.Warning("Could not unmarshal message from client");
				continue
			}

			switch wm.Level {
			case "debug":
				logrus.Debug(wm.Message)
			case "info":
				logrus.Info(wm.Message)
			case "warn":
				logrus.Warning(wm.Message)
			case "error":
				logrus.SetOutput(os.Stderr)
				logrus.Error(wm.Message)
				logrus.SetOutput(os.Stdout)
			default:
				logrus.Info(wm.Message)
			}
		}
	}
}