./toxiproxy-server 

./toxiproxy-cli create example.com --listen 0.0.0.0:8080 --upstream www.example.com:80
./toxiproxy-cli toxic add -t latency -a latency=1000 -u example.com
./toxiproxy-cli toxic add -t jitter -a jitter=900 -u example.com


## Change a toxic
./toxiproxy-cli toxic update -n latency_upstream -a latency=1000 -a jitter=900 example.com

## Inspect a proxy
./toxiproxy-cli toxic inspect example.com