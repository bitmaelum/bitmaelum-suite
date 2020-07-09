FROM debian:buster-slim

COPY ./release/linux/* /usr/bin/

EXPOSE 2424


RUN mkdir /etc/bitmaelum && ln -sf /bitmaelum/server-config.yml /etc/bitmaelum/server-config.yml

CMD ./usr/bin/bm-server
