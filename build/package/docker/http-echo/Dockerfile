FROM golang:alpine AS build-env

ARG REPONAME
ARG APPNAME
ENV SRCPATH $GOPATH/src/github.com/ewilde/$REPONAME
RUN apk update && apk add build-base
COPY ./ $SRCPATH
RUN go install github.com/ewilde/$REPONAME/cmd/$APPNAME

FROM alpine

ARG APPNAME
ARG TAGS
ENV BUILD_TAGS=$TAGS
ENV LOG_FORMAT=json
WORKDIR /app
COPY --from=build-env /go/bin/$APPNAME /app/
COPY build/package/docker/http-echo/wait-for /wait-for
EXPOSE 5678
ENTRYPOINT ["/wait-for"]
