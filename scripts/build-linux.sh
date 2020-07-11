#!/bin/sh

# This should be a Makefile I guess

REPO="github.com/bitmaelum/bitmaelum-suite"

APPS="bm-server bm-client bm-config bm-client-ui"
TOOLS="hash-address jwt proof-of-work readmail push-key"

TARGET=${1:-all}

export GOOS=linux
export GOARCH=amd64

# Generate commit / date variables we will inject in our code
BUILD_DATE=`date`
COMMIT=`git rev-parse HEAD`
PKG=`go list ./core`
GO_BUILD_FLAGS="-X '${PKG}.BuildDate=${BUILD_DATE}' -X '${PKG}.GitCommit=${COMMIT}'"

printf "Compiling ["

for APP in $APPS; do
  if [ "${TARGET}" = "all" ] || [ "${TARGET}" = "$APP" ] ; then
    go build -ldflags "${GO_BUILD_FLAGS}" -o release/linux/${APP} ${REPO}/cmd/${APP}
  fi
  printf "."
done

for TOOL in $TOOLS; do
  if [ "${TARGET}" = "all" ] || [ "${TARGET}" = "$TOOL" ] ; then
    go build -ldflags "${GO_BUILD_FLAGS}" -o release/linux/${TOOL} ${REPO}/tools/${TOOL}
  fi
  printf "."
done

printf "]\n"
