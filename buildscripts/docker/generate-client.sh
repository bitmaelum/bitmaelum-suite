#!/bin/sh

if [ ! -f Makefile ] ; then
  echo "Please run this file in the root of the source: ./buildscripts/docker/generate-client.sh"
  exit 1
fi

if [ ! -f ./release/linux-amd64/bm-client ] ; then
  make -j linux-amd64
fi

TMPDIR="tmp-$$"
mkdir $TMPDIR

cp ./buildscripts/docker/bitmaelum-client-config.yml $TMPDIR
cp ./buildscripts/docker/docker-entrypoint-client.sh $TMPDIR
cp ./release/linux-amd64/bm-* $TMPDIR

docker build -t bitmaelum/client -f ./buildscripts/docker/Dockerfile.client $TMPDIR
docker push bitmaelum/client:latest

if [ ! -z "$1" ] ; then
  echo "Tag found: $1"
  docker tag bitmaelum/client:latest bitmaelum/client:$1
  docker push bitmaelum/client:$1
fi


rm -rf $TMPDIR
