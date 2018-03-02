package main

import (
	"flag"
	"encoding/json"
	"os"
	"strings"

	"github.com/tdecker91/gws"
	log "github.com/sirupsen/logrus"
)

type websocketMessage struct {
	Level string `json:"level"`
	Message string `json:"message"`
}

var port = flag.Int("port", 8030, "port to listen on")
var route = flag.String("route", "/", "route to listen for socket connections")
var level = flag.String("level", "info", "logging level. [debug, info, warn, error]")
var formatter = flag.String("formatter", "text", "format to output log messages as. [text, json]")
var timeFormat = flag.String("timeFormat", "2006-01-02 15:04:05", "text formatter only. format string to output time logs. Format as defined here https://golang.org/pkg/time/#Time.Format")
var format = flag.String("format", "[%L] %D: %M", "text formatter only. format for the output message. %L = log level. %D = Date. %M = Message")

type textFormatter struct {
	LogrusFormatter *log.TextFormatter
}

func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	message := strings.Replace(*format, "%L", strings.ToUpper(entry.Level.String()), -1)
	message = strings.Replace(message, "%D", entry.Time.Format(*timeFormat), -1)
	message = strings.Replace(message, "%M", entry.Message, -1)

	return append([]byte(message), byte('\n')), nil
}

func init() {
	flag.Parse()

	switch *level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	switch *formatter {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		logrusFormatter := log.TextFormatter{}
		logrusFormatter.ForceColors = true
		logrusFormatter.FullTimestamp = true
		fmtr := &textFormatter{&logrusFormatter}
		log.SetFormatter(fmtr)
	}

}

func main() {
	messages := make(chan gws.Message)

	server := gws.NewSocketServer(*port, *route)
	go server.Start(messages)

	for {
		m := <-messages

		switch t := m.Type; t {
		case gws.ClientMessage:
			var wm websocketMessage
			if err := json.Unmarshal(m.Data, &wm); err != nil {
				log.Warning("Could not unmarshal message from client");
				continue
			}

			switch wm.Level {
			case "debug":
				log.Debug(wm.Message)
			case "info":
				log.Info(wm.Message)
			case "warn":
				log.Warning(wm.Message)
			case "error":
				log.SetOutput(os.Stderr)
				log.Error(wm.Message)
				log.SetOutput(os.Stdout)
			default:
				log.Info(wm.Message)
			}
		}
	}
}