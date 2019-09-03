FROM reg.example.com/base/alpine:3.6

MAINTAINER test

RUN apk update \
 && apk add --no-cache \
 && apk add ca-certificates

COPY configs/ /usr/local/etc/qcloud-callback-proxy/
COPY qcloud-callback-proxy /usr/local/bin/


EXPOSE 5050

CMD ["qcloud-callback-proxy", '', "-config=/usr/local/etc/qcloud-callback-proxy/proxy.toml"]
