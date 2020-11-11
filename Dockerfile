# -----------------------------------
FROM golang:1.15 AS builder

WORKDIR /app
COPY . /app
RUN make linux-amd64

# -----------------------------------
FROM debian:buster-slim

# We need CA certificates otherwise we cannot connect to https://
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netbase && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/release/linux-amd64/bm-* /usr/bin/
RUN mkdir /etc/bitmaelum && ln -sf /bitmaelum/server-config.yml /etc/bitmaelum/server-config.yml

EXPOSE 2424
CMD /usr/bin/bm-server
