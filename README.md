[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/EricNeid/go-relayserver) 

# go-relayserver

A simple websocket relay server, written in go. It relays incoming stream to multiple connected websockets.

## Install

```sh
 $ go get github.com/EricNeid/go-relayserver
```

## Usage

```sh
 $ go-relayserver optional: -port-stream <port> -port-ws <port> -s <secret>
```

## Testing

Make sure that ffmpeg is in your path. 
