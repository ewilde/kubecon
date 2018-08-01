# Evolving Systems Design: From Unreliable rpc to Resilience with Linkerd - Edward Wilde, Form3 (Intermediate Skill Level)                      
Project structure based on https://github.com/golang-standards/project-layout


# The demo
## Prerequisites
* [Docker](https://www.docker.com/community-edition#/download)
* [Go](https://golang.org/doc/install) and [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)

## Running the demo
```bash
export GOPATH=~/go  # or your alternative Go-language path
REPOPATH="$GOPATH"/src/github.com/ewilde/kubecon
mkdir -p "$REPOPATH"
git clone https://github.com/ewilde/kubecon.git "$REPOPATH"
cd "$REPOPATH"
make package
```

### Run system-1
`make up-s1`

### Run system-2
`make up-s2`



## Searching logs

### Service b - has delays sometime
docker logs system-2_service1b_1   2>&1  | grep Timeout

### Service c - error service
docker logs system-2_service1c_1   2>&1  | grep 503
