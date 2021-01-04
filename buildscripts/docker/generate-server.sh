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

cp ./buildscripts/docker/server-config.yml $TMPDIR
cp ./buildscripts/docker/docker-entrypoint-server.sh $TMPDIR
cp ./release/linux-amd64/bm-* $TMPDIR

docker build -t bitmaelum/server -f ./buildscripts/docker/Dockerfile.server $TMPDIR
docker push bitmaelum/server:latest

if [ ! -z "$1" ] ; then
  echo "Tag found: $1"
  docker tag bitmaelum/server:latest bitmaelum/server:$1
  docker push bitmaelum/server:$1
fi

rm -rf $TMPDIR
