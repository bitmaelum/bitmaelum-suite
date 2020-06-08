#!/bin/sh

#set -e

REPO="github.com/jaytaph/mailv2"

echo "[\c"

echo ".\c"
go build -o release/mailv2-server ${REPO}/server
echo ".\c"
go build -o release/proof-of-work ${REPO}/tools/proof-of-work
echo ".\c"
go build -o release/create-account ${REPO}/tools/create-account
echo ".\c"
go build -o release/protect-account ${REPO}/tools/protect-account
echo ".\c"
go build -o release/sendmail ${REPO}/tools/sendmail
echo ".\c"
go build -o release/client ${REPO}/client
echo ".\c"
go build -o release/jwt ${REPO}/tools/jwt

echo "]"
