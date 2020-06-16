#!/bin/sh

REPO="github.com/bitmaelum/bitmaelum-server"

TOOLS="create-account hash-address jwt proof-of-work protect-account readmail sendmail"

# We use govvv to inject GIT version information into the applications
go get github.com/ahmetb/govvv

GO_PATH=`go env GOPATH`
GO_BUILD_FLAGS=`${GO_PATH}/bin/govvv build -pkg version -flags`

echo "Compiling [\c"

echo ".\c"
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-server ${REPO}/bm-server
echo ".\c"
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-config ${REPO}/bm-config
echo ".\c"
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-client ${REPO}/bm-client-ui
echo ".\c"

for TOOL in $TOOLS; do
  go build -ldflags "${GO_BUILD_FLAGS}" -o release/${TOOL} ${REPO}/tools/${TOOL}
  echo ".\c"
done

echo "]"
