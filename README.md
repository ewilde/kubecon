# Evolving Systems Design: From Unreliable rpc to Resilience with Linkerd - Edward Wilde, Form3 (Intermediate Skill Level)                      

Project structure based on https://github.com/golang-standards/project-layout


# Demo
## Searching logs

### Service b - has delays sometime
docker logs system2_service1b_1   2>&1  | grep Timeout

### Service c - error service
docker logs system2_service1c_1   2>&1  | grep 503


### Load-balancer errors

# Toxi-proxy research

./toxiproxy-server 

./toxiproxy-cli create example.com --listen 0.0.0.0:8080 --upstream www.example.com:80
./toxiproxy-cli toxic add -t latency -a latency=1000 -u example.com
./toxiproxy-cli toxic add -t jitter -a jitter=900 -u example.com


## Change a toxic
./toxiproxy-cli toxic update -n latency_upstream -a latency=1000 -a jitter=900 example.com

## Inspect a proxy
./toxiproxy-cli toxic inspect example.com
