#!/bin/sh

#
# Create a self-signed certificate for the mailserver.
#
# Usage: ./create-cert.sh <my.server.tld>
#

openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 3650 -nodes -subj "/CN=${1}"
