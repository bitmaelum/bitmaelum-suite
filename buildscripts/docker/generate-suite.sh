#!/bin/sh

if [ ! -f Makefile ] ; then
  echo "Please run this file in the root of the source: ./buildscripts/docker/generate-server.sh"
  exit 1
fi

if [ ! -f ./release/linux-amd64/bm-server ] ; then
  make -j linux-amd64
fi

TMPDIR="tmp-$$"
mkdir $TMPDIR

#cp ./buildscripts/docker/server-config.yml $TMPDIR
cp ./buildscripts/docker/docker-entrypoint-suite.sh $TMPDIR
cp ./release/linux-amd64/bm-* $TMPDIR

docker build -t bitmaelum/suite -f ./buildscripts/docker/Dockerfile.suite $TMPDIR
docker push bitmaelum/suite:latest

if [ ! -z "$1" ] ; then
  echo "Tag found: $1"
  docker tag bitmaelum/suite:latest bitmaelum/suite:$1
  docker push bitmaelum/suite:$1
fi

rm -rf $TMPDIR
