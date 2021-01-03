FROM debian:buster-slim

# We need CA certificates otherwise we cannot connect to https://
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netbase && rm -rf /var/lib/apt/lists/*

COPY bm-* /usr/local/bin/
COPY client-config.yml /etc/bitmaelum/client-config.yml

COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
