package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type event struct {
	id        string
	username  string
	timestamp string
	port      string
	success   string
	ipaddress string
	exec      string
	cmd       string
}

func (e *event) readLine(message string) error {
	words := strings.Fields(message)
	for _, word := range words {
		e.readWord(word)
	}
	return nil
}

func (e *event) readWord(s string) {
	if strings.Contains(s, equals) {
		outer := strings.Split(s, equals)
		s = strings.Replace(s, prefix, "", 1)
		switch outer[0] {
		case "type":
			if outer[1] == "EOE" {
				e.flush()
			}
		case "msg":
			if strings.Contains(s, colon) {
				inner := strings.Split(s, colon)
				e.id = strings.TrimRight(inner[1], rightParen)
				e.timestamp = strings.Split(inner[0], dot)[0]
			}
		case "uid":
			e.username = outer[1]
		case "success":
			e.success = outer[1]
		case "exe":
			e.exec = outer[1]
		case "comm":
			e.cmd = outer[1]
		case "saddr":
			saddr, _ := hex.DecodeString(strings.TrimLeft(outer[1], "saddr="))
			port, _ := strconv.ParseInt(fmt.Sprintf("%x", saddr[2:4]), 16, 0)
			e.port = fmt.Sprintf("%v", port)
			if strings.HasPrefix(outer[1], ipv4Addr) {
				e.ipaddress = net.IP(saddr[4:8]).String()
			}
			if strings.HasPrefix(outer[1], ipv6Addr) {
				e.ipaddress = net.IP(saddr[8:24]).String()
			}
		}
	}
}

func (e *event) flush() *event {
	if e.username != "" && e.port != "" && e.ipaddress != "" {
		log.Println(e)
	}
	e.id = ""
	e.username = ""
	e.timestamp = ""
	e.port = ""
	e.success = ""
	e.ipaddress = ""
	e.exec = ""
	e.cmd = ""
	return e
}

func (e event) String() string {
	return (toUnixTimestamp(e.timestamp) + " uid: " + e.username + " destination: " + e.ipaddress + " port: " + e.port + " command: " + e.cmd + " exec: " + e.exec + " success: " + e.success)
}
