package main

import "fmt"
import "strings"
import "time"
import "strconv"

import "os"
import "log"
import "encoding/hex"
import "net"
import "gopkg.in/natefinch/lumberjack.v2"
import "bufio"

var rule auditRule
var event auditEvent

type auditRule struct {
	name  string
	key   string
	value string
}

type auditEvent struct {
	username  string
	timestamp string
	port      string
	success   string
	ipaddress string
	exec      string
	cmd       string
}

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/audit/audit-dispatcher.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     7,
	})

	log.SetFlags(0) // turning off timestamps
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parseIt(scanner.Text())
	}
}

func flushEvent() {
	if event.username == "" {
		return
	}
	if event.username == "0" {
		return
	}
	if event.port == "" {
		return
	}
	if event.ipaddress == "" {
		return
	}
	if event.ipaddress == "127.0.0.1" { // probably want to short circut these way before the purge event
		return
	}
	if event.ipaddress == "0.0.0.0" {
		return
	}
	if strings.HasPrefix(event.ipaddress, "199.116.7") { // close enough
		return
	}
	if strings.HasPrefix(event.ipaddress, "192.168") { // close enough
		return
	}
	if strings.HasPrefix(event.ipaddress, "104.37.8") {
		return
	}
	//   if (event.success == "no") {
	//      return
	//   }
	log.Println(toUnixTimestamp(event.timestamp) + " uid: " + event.username + " destination: " + event.ipaddress + " port: " + event.port + " command: " + event.cmd + " exec: " + event.exec + " success: " + event.success)
	event.purgeEvent()
	return
}

func parseIt(message string) {
	words := strings.Fields(message)
	for _, word := range words {
		if strings.HasPrefix(word, "msg=audit(") { // we found the timestamp
			word = strings.TrimRight(strings.TrimLeft(word, "msg=audit("), "):") // trim for unixtime timestamp
			event.setTs(strings.Split(strings.TrimRight(strings.TrimLeft(word, "msg=audit("), "):"), ".")[0])
			if rule.name == word {
				flushEvent()
			}
			rule.setName(word)
		}
		if strings.HasPrefix(word, "uid=") { // we found a uid
			//       usr, err := user.LookupId(strings.TrimLeft(word, "uid="))
			//       if err != nil {
			//         event.setUsername(strings.TrimLeft(word, "uid="))
			//       } else {
			event.setUsername(strings.TrimLeft(word, "uid="))
			//       }
			continue
		}
		if strings.HasPrefix(word, "success=") {
			event.setSuccess(strings.TrimLeft(word, "success="))
			continue
		}
		if strings.HasPrefix(word, "exe=") {
			event.setExec(strings.TrimLeft(word, "exe="))
			continue
		}
		if strings.HasPrefix(word, "comm=") {
			event.setCmd(strings.TrimLeft(word, "comm="))
			continue
		}
		if strings.HasPrefix(word, "saddr=0200") { // != 0200 is the wrong type of socket address family
			saddr, _ := hex.DecodeString(strings.TrimLeft(word, "saddr="))
			result, _ := strconv.ParseInt(fmt.Sprintf("%x", saddr[2:4]), 16, 0)
			event.setPort(fmt.Sprintf("%v", result))
			event.setIpAddress(net.IP(saddr[4:8]).String())
			continue
		}
	}
}

func (rule auditRule) getKey(input string) string {
	s := strings.Split(input, "=")
	key := s[0]
	return key
}

func (rule auditRule) getValue(input string) string {
	s := strings.Split(input, "=")
	value := s[1]
	return value
}
func (event *auditEvent) setTs(input string) *auditEvent {
	event.timestamp = input
	return event
}
func (event *auditEvent) setUsername(input string) *auditEvent {
	event.username = input
	return event
}
func (event *auditEvent) setSuccess(input string) *auditEvent {
	event.success = input
	return event
}
func (event *auditEvent) setIpAddress(input string) *auditEvent {
	event.ipaddress = input
	return event
}

func (event *auditEvent) purgeEvent() *auditEvent {
	event.username = ""
	event.timestamp = ""
	event.port = ""
	event.success = ""
	event.ipaddress = ""
	return event
}
func (event *auditEvent) setCmd(input string) *auditEvent {
	event.cmd = input
	return event
}
func (event *auditEvent) setExec(input string) *auditEvent {
	event.exec = input
	return event
}
func (event *auditEvent) setPort(input string) *auditEvent {
	if input != "0" {
		event.port = input
	}
	return event
}

func (rule *auditRule) setName(input string) *auditRule {
	rule.name = input
	return rule
}

func (rule *auditRule) setKey(input string) *auditRule {
	rule.key = input
	return rule
}

func (rule *auditRule) setVaue(input string) *auditRule {
	rule.value = input
	return rule
}

func toUnixTimestamp(input string) string {
	// time.Unix(1392899576, 0).Format(time.RFC3339)
	if s, err := strconv.ParseInt(input, 10, 64); err == nil {
		return time.Unix(s, 0).Format(time.RFC3339)
	}
	return ""
}
