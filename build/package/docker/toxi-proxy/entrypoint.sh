#!/usr/local/bin/dumb-init /bin/bash
set -ex

exec "$@"  >& _toxi.out

./toxiproxy-cli create example.com --listen 0.0.0.0:8080 --upstream www.example.com:80
./toxiproxy-cli toxic add -t latency -a latency=1000 -u example.com
./toxiproxy-cli toxic add -t jitter -a jitter=900 -u example.com


tail -f _toxi.out