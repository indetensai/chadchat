#!/bin/sh
if [ ! -f access_private.pem ]; then
    echo "File access_private.pem does not exist. Generating."
    openssl genrsa 2048 | openssl pkey -traditional -out access_private.pem
fi

if [ ! -f refresh_private.pem ]; then
    echo "File refresh_private.pem does not exist. Generating."
    openssl genrsa 2048 | openssl pkey -traditional -out refresh_private.pem
fi