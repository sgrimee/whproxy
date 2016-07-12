# webhookProxy

Proxy incoming webhooks to established websockets.

This proxy can be used when an application needs to receive webhooks from some platform and a firewall blocks incoming connections. When client applications establish a websocket connection to this proxy, a listener is created for incoming webhooks. When a webhook arrives on a given listener, the event is signaled to the corresponding client app via the websocket.

## Status

Work in progress, not functional
