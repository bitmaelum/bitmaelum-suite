#!/bin/sh

REPO="github.com/jaytaph/mailv2"

go build -o release/mailv2-server ${REPO}/app/server
go build -o release/mailv2-client ${REPO}/app/client
go build -o release/proof-of-work ${REPO}/app/proof-of-work
