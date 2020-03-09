# Syslog_converter

![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/majestik/syslog_converter)

## Purpose
Converts RFC 3164 UDP syslog messages to RFC 5424 TCP messages. 
Intended to convert messages for Promtail ingest

## Usage

Docker:
 
 `docker run -p 1514:514/udp syslog_converter:latest listen -t 172.16.14.244 -p 1514 -l 514`

Comnand line:

`syslog_converter listen -t 172.16.14.244 -p 1514 -l 514`


## Build

Docker:

`docker build -t syslog_converter .`

Development:

`go get -d -v ; go build `
