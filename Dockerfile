FROM alpine:3.4
MAINTAINER Sam Grimee <sgrimee gmail.com>
ADD whproxy_linux /whproxy
ENTRYPOINT /whproxy