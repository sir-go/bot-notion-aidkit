package main

import (
	"log"
	"os"
	"os/signal"
)

var (
	CFG *Config
	LOG *log.Logger
)

func initInterrupt() {
	LOG.Println("-- start --")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(c chan os.Signal) {
		for range c {
			LOG.Println("-- stop --")
			os.Exit(137)
		}
	}(c)
}

func init() {
	CFG = ConfigInit()
	LOG = initLogging()
	initInterrupt()
}
