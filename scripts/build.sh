#!/bin/sh

REPO="github.com/bitmaelum/bitmaelum-server"

TOOLS="create-account hash-address jwt proof-of-work protect-account readmail sendmail"

# Generate commit / date variables we will inject in our code
BUILD_DATE=`date`
COMMIT=`git rev-parse HEAD`
PKG=`go list ./core`
GO_BUILD_FLAGS="-X '${PKG}.BuildDate=${BUILD_DATE}' -X '${PKG}.GitCommit=${COMMIT}'"

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
