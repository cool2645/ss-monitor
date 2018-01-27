FROM alpine
RUN apk add --update curl \
    python \
    py-pip \
    git \
    proxychains-ng \
    libsodium \
    && pip install git+https://github.com/shadowsocks/shadowsocks.git@master
COPY proxychains.conf /etc/proxychains/proxychains.conf
COPY test-your-ss.sh /
COPY shadowsocks.json /
CMD ["/bin/sh", "test-your-ss.sh"]