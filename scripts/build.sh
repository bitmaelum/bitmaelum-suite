#!/bin/sh

# This should be a Makefile I guess

REPO="github.com/bitmaelum/bitmaelum-server"

TOOLS="create-account hash-address jwt proof-of-work protect-account readmail sendmail"

TARGET=${1:-all}
# Generate commit / date variables we will inject in our code
BUILD_DATE=`date`
COMMIT=`git rev-parse HEAD`
PKG=`go list ./core`
GO_BUILD_FLAGS="-X '${PKG}.BuildDate=${BUILD_DATE}' -X '${PKG}.GitCommit=${COMMIT}'"

printf "Compiling ["

printf  "."
if [[ ${TARGET} == "all" || ${TARGET} == "bm-server" ]] ; then
  go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-server ${REPO}/bm-server
fi

printf  "."
if [[ ${TARGET} == "all" || ${TARGET} == "bm-config" ]] ; then
  go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-config ${REPO}/bm-config
fi

printf  "."
if [[ ${TARGET} == "all" || ${TARGET} == "bm-client-ui" ]] ; then
  go build -ldflags "${GO_BUILD_FLAGS}" -o release/bm-client ${REPO}/bm-client-ui
fi
printf  "."

for TOOL in $TOOLS; do
  if [[ ${TARGET} == "all" || ${TARGET} == $TOOL ]] ; then
    go build -ldflags "${GO_BUILD_FLAGS}" -o release/${TOOL} ${REPO}/tools/${TOOL}
  fi
  printf "."
done

printf "]\n"
