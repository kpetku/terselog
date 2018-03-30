package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	if len(os.Args[1:]) == 1 {
		log.SetOutput(loadConfig(os.Args[1]))
	} else {
		l := lumberjack.Logger{}
		l.Filename = "/var/log/audit/terselog.log"
		l.MaxSize = 10
		l.MaxBackups = 10
		l.MaxAge = 7
		log.SetOutput(&l)
	}
	log.SetFlags(0) // turning off timestamps

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		for {
			<-c
			log.Println("Exiting because auditd is stopping")
			os.Exit(0)
		}
	}()
	scanner := bufio.NewScanner(os.Stdin)
	var e event
	for scanner.Scan() {
		e.readLine(scanner.Text())
	}
}
