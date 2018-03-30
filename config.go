package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func loadConfig(s string) *lumberjack.Logger {
	l := lumberjack.Logger{}
	f, err := os.Open(s)
	if err != nil {
		log.Fatalf("%s", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		split := strings.Fields(scanner.Text())
		if len(split) != 2 {
			log.Fatalf("Malformed line: %s: too short", split)
		}
		switch split[0] {
		case "Filename":
			l.Filename = split[1]
		case "MaxSize":
			l.MaxAge, err = convert(split[1])
		case "MaxBackups":
			l.MaxBackups, err = convert(split[1])
		case "MaxAge":
			l.MaxAge, err = convert(split[1])
		}
		if err != nil {
			log.Fatalf("Malformed line: %s: %s", split, err)
		}
	}
	defer f.Close()
	return &l
}
