# terselog
Timestamped outgoing TCP IPv4 and IPv6 connection logs *for humans* via the auditd subsystem.

## Example output
```
2018-03-29T21:16:46-04:00 uid: 1000 destination: 54.244.19.239 port: 0 command: 444E53205265737E65722023323839 exec: "/usr/lib/firefox/firefox" success: yes
2018-03-29T21:16:46-04:00 uid: 1000 destination: 127.0.0.53 port: 53 command: "curl" exec: "/usr/bin/curl" success: yes
2018-03-29T21:16:46-04:00 uid: 1000 destination: 2a03:2880:f127:283:face:b00c:0:25de port: 80 command: "curl" exec: "/usr/bin/curl" success: yes
2018-03-29T21:16:46-04:00 uid: 1000 destination: 157.240.2.35 port: 80 command: "curl" exec: "/usr/bin/curl" success: yes
```

## Install on Ubuntu 17.10
### Install the required dependencies
```
sudo apt-get -y install auditd audispd-plugins
```
### Copy the `terselog` binary to `/sbin/terselog`
```
cp terselog /sbin/terselog
```
### Create a file named `/etc/terselog.conf` containing
```
Filename /var/log/audit/terselog.log
MaxSize 10
MaxBackups 10
MaxAge 7
```
#### The following options are taken from lumberjack's [documentation](https://godoc.org/github.com/natefinch/lumberjack#Logger):
- **Filename** Filename is the file to write logs to.  Backup log files will be retained in the same directory.
- **MaxSize** is the maximum size in megabytes of the log file before it gets rotated.
- **MaxAge** is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
- **MaxBackups** is the maximum number of old log files to retain.  The default is to retain all old log files (though MaxAge may still cause them to get deleted.)

### Create a file named `/etc/audisp/plugins.d/terselog.conf` containing
```
active = yes
direction = out
path = /sbin/terselog
type = always
format = string
```
### Create a file named `/etc/audit/rules.d/terselog.rules` containing
```
-a exit,always -F arch=b64 -S connect -F a2!=110 -k outbound
```
### Ensure the `terselog` binary is root owned otherwise audisp *will not* execute and terselog will fail to run
```
chown root:root /sbin/terselog
```
### Restart auditd:
```
systemctl restart auditd.service
```
## Dependencies
lumberjack by natefinch, auditd, audisp, and Go.