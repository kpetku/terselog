package main

import (
	"strconv"
	"time"
)

const (
	prefix     string = "msg=audit("
	equals     string = "="
	colon      string = ":"
	dot        string = "."
	rightParen string = ")"
	ipv4Addr   string = "02"
	ipv6Addr   string = "0A"
)

func toUnixTimestamp(s string) string {
	if t, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(t, 0).Format(time.RFC3339)
	}
	return ""
}

func convert(s string) (int, error) {
	return strconv.Atoi(s)
}
