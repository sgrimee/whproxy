# whproxy - Web Hooks Proxy

Proxy incoming web hooks to established web sockets.

This proxy can be used when an application needs to receive json web hooks from some platform and a firewall blocks incoming connections. When client applications establish a web socket connection to this proxy, a listener is created for incoming web hooks. When a web hook arrives on a given listener, the event is signaled to the corresponding client application via the web socket.

## Usage

whproxy --host my.publichost.com --port 12345

At the moment the same port is used for webhooks and websockets.

Then establish a websocket connection to my.publichost.com:12345 and you will receive a json payload with the URL of your private webhook endpoint.

## Status

Work in progress, functional, barely tested.

