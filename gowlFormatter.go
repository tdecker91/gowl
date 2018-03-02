package main

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type GowlFormatter struct {
	FormatString string
}

func (gw *GowlFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	message := strings.Replace(gw.FormatString, "%L", strings.ToUpper(entry.Level.String()), -1)
	message = strings.Replace(message, "%D", entry.Time.Format(*timeFormat), -1)
	message = strings.Replace(message, "%M", entry.Message, -1)

	return append([]byte(message), byte('\n')), nil
}