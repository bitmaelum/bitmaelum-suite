# -----------------------------------
FROM golang:1.15 AS builder

WORKDIR /app
COPY . /app
RUN make linux-amd64

# -----------------------------------
FROM debian:buster-slim

COPY --from=builder /app/release/linux-amd64/* /usr/bin/
RUN mkdir /etc/bitmaelum && ln -sf /bitmaelum/server-config.yml /etc/bitmaelum/server-config.yml

EXPOSE 2424
CMD ./usr/bin/bm-server
