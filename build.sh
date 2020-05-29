#!/bin/sh

set -e

REPO="github.com/jaytaph/mailv2"

go build -o release/mailv2-server ${REPO}/server
go build -o release/proof-of-work ${REPO}/tools/proof-of-work
go build -o release/create-account ${REPO}/tools/create-account
go build -o release/sendmail ${REPO}/tools/sendmail
