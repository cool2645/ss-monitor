FROM alpine
RUN apk add --update bash curl
COPY hosts /etc/hosts
COPY introspect.sh /
COPY beacon.txt /
CMD ["/bin/bash", "introspect.sh"]
