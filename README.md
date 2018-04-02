# Evolving Systems Design: From Unreliable rpc to Resilience with Linkerd - Edward Wilde, Form3 (Intermediate Skill Level)                      
Project structure based on https://github.com/golang-standards/project-layout


# The demo
## Prerequisites
* [Docker](https://www.docker.com/community-edition#/download)

## Running the demo
```bash
git@github.com:ewilde/kubecon.git ewilde-kubecon
cd ewilde-kubecon
make package
```

### Run system-1
`make up-s1`

### Run system-2
`make up-s2`



## Searching logs

### Service b - has delays sometime
docker logs system2_service1b_1   2>&1  | grep Timeout

### Service c - error service
docker logs system2_service1c_1   2>&1  | grep 503
