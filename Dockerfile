FROM debian:buster-slim

RUN mkdir /etc/bitmaelum && ln -sf /bitmaelum/server-config.yml /etc/bitmaelum/server-config.yml
COPY ./release/linux-amd64/* /usr/bin/

EXPOSE 2424
CMD ./usr/bin/bm-server
