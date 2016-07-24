# whproxy - Web Hooks Proxy

Proxy incoming web hooks to established web sockets.

This proxy can be used when an application needs to receive web hooks from some platform and a firewall blocks incoming connections. When client applications establish a web socket connection to this proxy, a listener is created for incoming web hooks. When a web hook arrives on a given listener, the event is signaled to the corresponding client application via the web socket.

## Status

Work in progress, not functional
