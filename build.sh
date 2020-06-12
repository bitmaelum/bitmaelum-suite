#!/bin/sh

#set -e

REPO="github.com/jaytaph/mailv2"

TOOLS="create-account hash-address jwt mail-server-config proof-of-work protect-account readmail sendmail"

echo "Compiling [\c"

echo ".\c"
go build -o release/mailv2-server ${REPO}/server
echo ".\c"
go build -o release/client ${REPO}/client-ui
echo ".\c"

for TOOL in $TOOLS; do
  go build -o release/${TOOL} ${REPO}/tools/${TOOL}
  echo ".\c"
done

echo "]"
