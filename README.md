# whproxy - Web Hooks Proxy

Proxy incoming web hooks to established web sockets.

This proxy can be used when an application needs to receive json web hooks from some platform and a firewall blocks incoming connections. When client applications establish a web socket connection to this proxy, a listener is created for incoming web hooks. When a web hook arrives on a given listener, the event is signaled to the corresponding client application via the web socket.

## Usage

    whproxy --host my.publichost.com --port 8080

At the moment the same port is used for webhooks and websockets.

Then establish a websocket connection to my.publichost.com:8080 and you will receive a json payload with the URL of your private webhook endpoint.
Incoming webhooks events are passed on to the websocket, encapsulated inside the 'data' field of the json object.

## Encryption

Encryption is enabled by providing a certificate and key file

    whproxy --host my.publichost.com --port 8443 --cert cert.pem --key key.unencrypted.pem

If you use self-signed certificates, be sure to enable cert pinning on the client.

    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 360
    openssl rsa -in key.pem -out key.unencrypted.pem

## Validation

You can enforce the presence of a valid HMAC in the header of incoming webhooks but giving the '-validate' option

## Docker

    docker run -it -p 12345:12345 quay.io/sgrimee/whproxy:v0.2.0 -host my.docker.host

Available versions can be found at: https://quay.io/repository/sgrimee/whproxy?tab=tags

If you use certificates, mount them with

```
    docker run -it -p 8080:12345 \
      -v $(pwd)/cert.pem:/cert.pem \
      -v $(pwd)/key.unencrypted.pem:/key.pem \
      quay.io/sgrimee/whproxy:v0.2.0 \
        -host my.docker.host \
        -cert /cert.pem -key /key.pem
```

