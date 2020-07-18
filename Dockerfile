FROM debian:buster-slim

# Needed to connect to https sites
RUN apt-get update && apt-get install ca-certificates -y && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN mkdir /etc/bitmaelum && ln -sf /bitmaelum/server-config.yml /etc/bitmaelum/server-config.yml
COPY ./release/linux-amd64/* /usr/bin/

EXPOSE 2424
CMD ./usr/bin/bm-server
