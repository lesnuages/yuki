# Yuki

## Get & Run

Set your `$GOPATH` environment variable and just run:

```bash
$ go get -u github.com/lesnuages/yuki
```

Then, you can start with:

```bash
$ go run yuki.go -f <pcap_filepath>
```

## Description

A rewrite of [clipp](https://github.com/lesnuages/clipp) in Go. Supports both PCAP and PCAPNG file formats thanks to [GoPacket](https://github.com/google/gopacket)

## Features

The following features are (or are planned to be) implemented:

* Session listing
* Session export to PCAP format
* Full text search
* Application layer export to binary format
* ...
