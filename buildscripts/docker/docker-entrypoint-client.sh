#!/bin/sh

# Check if the /bitmaelum directory is mounted. If not, throw an error
mount | grep 'on /bitmaelum type'
if [ $? -eq 1 ] ; then
  echo ""
  echo "********************************* WARNING *********************************"
  echo ""
  echo "You are trying to run the bitmaelum docker image without mounting a local "
  echo "directory to /bitmaelum. This means your data will not be persisted and "
  echo "you could loose your vault data. Please mount the directory through docker:"
  echo ""
  echo "     docker run -v /local/dir:/bitmaelum bitmaelum/client:latest"
  echo ""
  echo "********************************* WARNING *********************************"
  exit 0
fi

exec /usr/local/bin/bm-client "$@"
