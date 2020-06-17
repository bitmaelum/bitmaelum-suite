#!/bin/sh

REPO="github.com/bitmaelum/bitmaelum-server"

TOOLS="create-account hash-address jwt proof-of-work protect-account readmail sendmail"

# We use govvv to inject GIT version information into the applications
go get github.com/ahmetb/govvv

GO_PATH=`go env GOPATH`
GO_BUILD_FLAGS=`${GO_PATH}/bin/govvv build -pkg version -flags`

printf "Compiling ["

printf  "."
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-server ${REPO}/bm-server
printf  "."
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-config ${REPO}/bm-config
printf  "."
go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-client ${REPO}/bm-client-ui
printf  "."

for TOOL in $TOOLS; do
  go build -ldflags "${GO_BUILD_FLAGS}" -o release/${TOOL} ${REPO}/tools/${TOOL}
  printf "."
done

printf "]\n"
