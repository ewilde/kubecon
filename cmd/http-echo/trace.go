package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cenkalti/backoff"
	"github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/pkg/errors"
)

var collector zipkintracer.Collector

var once sync.Once
var zipkinHostAndPort = fmt.Sprintf("%s:%s",
	valueOrDefault(os.Getenv("ZIPKIN_HOST"), "localhost"),
	valueOrDefault(os.Getenv("ZIPKIN_PORT"), "9410"))

func trace(request *http.Request, logOut io.Writer, serviceName string, message string, duration time.Duration) error {
	once.Do(func() {
		var err error
		collector, err = zipkintracer.NewScribeCollector(zipkinHostAndPort, time.Second, zipkintracer.ScribeBatchSize(0), zipkintracer.ScribeBatchInterval(time.Millisecond))
		if err != nil {
			log.Fatal(err)
		}
	})

	var spanID = int64(rangeIn(100000, 999999))
	parentID, _, traceID, _, err := getParentSpan(logOut, request)
	if err != nil {
		return err
	}

	if traceID == 0 {
		traceID = int64(rangeIn(100000, 999999))
	}

	timeStamp := time.Now().Add(duration * -1)
	annotation := &zipkincore.Annotation{
		Timestamp: timeStamp.UnixNano() / 1e3,
		Value:     "cs",
		Host: &zipkincore.Endpoint{
			ServiceName: serviceName,
			Ipv4:        int32(binary.BigEndian.Uint32(ipAddress)),
		},
	}

	span := makeNewSpan(timeStamp, message, traceID, spanID, parentID, []*zipkincore.Annotation{annotation}, duration, true)

	fmt.Fprint(logOut, fmt.Sprintf("[DEBUG] Traceid:%d spanid:%d parentid:%d duration:%v ms\n", traceID, spanID, parentID, duration.Seconds()*1000))
	return collector.Collect(span)
}

func getParentSpan(logOut io.Writer, request *http.Request) (spanID int64, parentID int64, traceID int64, flags int64, err error) {
	traceHeaders, ok := request.Header["L5d-Ctx-Trace"]
	if !ok {
		return 0, 0, 0, 0, nil
	}

	traceBytes, err := base64.StdEncoding.DecodeString(traceHeaders[0])
	if err != nil {
		fmt.Fprint(logOut, err)
		return 0, 0, 0, 0, err
	}

	if len(traceBytes) != 32 {
		fmt.Fprint(logOut, "[WARN] Expected 32 bytes")

		return 0, 0, 0, 0, errors.New(fmt.Sprintf("Expected 32 bytes, got %d", len(traceBytes)))
	}

	spanID = int64(binary.BigEndian.Uint64(traceBytes[:8]))
	parentID = int64(binary.BigEndian.Uint64(traceBytes[8:16]))
	traceID = int64(binary.BigEndian.Uint64(traceBytes[16:24]))
	flags = int64(binary.BigEndian.Uint64(traceBytes[24:32]))
	return spanID, parentID, traceID, flags, nil
}

func makeNewSpan(startTime time.Time, methodName string, traceID, spanID, parentSpanID int64, annotations []*zipkincore.Annotation, duration time.Duration, debug bool) *zipkincore.Span {
	timestamp := startTime.UnixNano() / 1e3
	return &zipkincore.Span{
		TraceID:     traceID,
		Name:        methodName,
		ID:          spanID,
		ParentID:    &parentSpanID,
		Debug:       debug,
		Duration:    thrift.Int64Ptr(duration.Nanoseconds() / 1000),
		Timestamp:   &timestamp,
		Annotations: annotations,
	}
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func retry(op func() error, timeout time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = timeout
	return backoff.Retry(op, bo)
}
