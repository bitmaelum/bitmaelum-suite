#!/bin/sh

if [ ! -f Makefile ] ; then
  echo "Please run this file in the root of the source: ./buildscripts/docker/generate-docker.sh"
  exit 1
fi

if [ ! -f ./release/linux-amd64/bm-client ] ; then
  make -j linux-amd64
fi

TMPDIR="tmp-$$"
mkdir $TMPDIR

cp ./buildscripts/docker/client-config.yml $TMPDIR
cp ./buildscripts/docker/docker-entrypoint.sh $TMPDIR
cp ./release/linux-amd64/bm-* $TMPDIR

docker build -t bitmaelum/bitmaelum-client -f ./buildscripts/docker/Dockerfile $TMPDIR
docker push bitmaelum/bitmaelum-client:latest

rm -rf $TMPDIR
