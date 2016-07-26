FROM scratch
MAINTAINER Sam Grimee <sgrimee gmail.com>
ADD whproxy_linux /whproxy
ENTRYPOINT ["/whproxy"]