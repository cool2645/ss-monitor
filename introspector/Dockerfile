FROM alpine
RUN apk add --update curl
COPY hosts /etc/hosts
COPY introspect.sh /
COPY beacon.txt /
CMD ["/bin/sh", "introspect.sh"]
