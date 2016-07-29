# whproxy - Web Hooks Proxy

Proxy incoming web hooks to established web sockets.

This proxy can be used when an application needs to receive json web hooks from some platform and a firewall blocks incoming connections. When client applications establish a web socket connection to this proxy, a listener is created for incoming web hooks. When a web hook arrives on a given listener, the event is signaled to the corresponding client application via the web socket.

## Usage

    whproxy --host my.publichost.com --port 12345

At the moment the same port is used for webhooks and websockets.

Then establish a websocket connection to my.publichost.com:12345 and you will receive a json payload with the URL of your private webhook endpoint.
Incoming webhooks events are passed on to the websocket, encapsulated inside the 'data' field of the json object.

## Encryption

Encryption is enabled by providing a certificate and key file

    whproxy --host my.publichost.com --port 80 --sslport 443 --cert cert.pem --key key.unencrypted.pem

If you use self-signed certificates, be sure to enable cert pinning on the client.

    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 360
    openssl rsa -in key.pem -out key.unencrypted.pem -passin pass:TYPE_YOUR_PASS

## Docker

    docker run -it -p 12345:12345 quay.io/sgrimee/whproxy -host my.docker.host

If you use certificates, mount them with

    docker run -it -p 12345:12345 -p 12346:12346 \
      -v cert.pem:/cert.pem -v key.unencrypted.pem:/key.pem \
      quay.io/sgrimee/whproxy -host my.docker.host -cert /cert.pem -key /key.pem


