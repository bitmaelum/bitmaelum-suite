#!/bin/sh

#set -e

REPO="github.com/bitmaelum/bitmaelum-server"

TOOLS="create-account hash-address jwt mail-server-config proof-of-work protect-account readmail sendmail"

GO_PATH=`go env GOPATH`
GO_BUILD_FLAGS=`${GO_PATH}/bin/govvv build -pkg version -flags`

echo "Compiling [\c"

echo ".\c"
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bitmaelum-server ${REPO}/server
echo ".\c"
go build -ldflags "${GO_BUILD_FLAGS}" -o release/client ${REPO}/client-ui
echo ".\c"

for TOOL in $TOOLS; do
  go build -ldflags "${GO_BUILD_FLAGS}" -o release/${TOOL} ${REPO}/tools/${TOOL}
  echo ".\c"
done

echo "]"
