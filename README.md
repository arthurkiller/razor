# Razor
tcp-fast-open supported tcp connection utility in go

make sure your server suppport tcp-fast-open and enabled

## Check your system
* Kernal should above 3.7.0

> uname -a

* Enable the TFO, /proc/sys/net/ipv4/tcp_fastopen should be 3

> cat /proc/sys/net/ipv4/tcp_fastopen

## Feature
* Tested on the CentOS 6 with kernal 3.10.*
* Easy to use. Implement net.Conn
* Support the TCP Fast Open

## TODO
* Cross platform. Linux (MacOS & Windows ___TODO___)
* Fix the write for write error

## Contribute
Wellcome any PR, just put up an issue!
