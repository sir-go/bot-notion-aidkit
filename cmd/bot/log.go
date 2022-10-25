package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s %s", time.Now().Local().Format("06/01/02 15:04:05"), string(bytes))
}

func initLogging() (l *log.Logger) {
	l = log.New(os.Stdout, "", log.Lshortfile)
	l.SetOutput(new(logWriter))
	return
}
