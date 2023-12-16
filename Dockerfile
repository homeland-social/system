FROM alpine:3.18

ADD system /system

ENTRYPOINT ["/system"]
