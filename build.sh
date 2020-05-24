#!/bin/sh

REPO="github.com/jaytaph/mailv2"

go build -o release/mailv2-server ${REPO}/server
go build -o release/mailv2-client ${REPO}/client
go build -o release/proof-of-work ${REPO}/tools/proof-of-work
