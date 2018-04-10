package main

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var collector zipkintracer.Collector

var traceID = int64(rangeIn(100000, 999999))
var spanID = int64(rangeIn(100000, 999999))
var once sync.Once
var zipkinHostAndPort = fmt.Sprintf("%s:%s",
	valueOrDefault(os.Getenv("ZIPKIN_HOST"), "localhost"),
	valueOrDefault(os.Getenv("ZIPKIN_PORT"), "9410"))

func trace(logOut io.Writer, message string, duration time.Duration) error {
	once.Do(func() {
		var err error
		collector, err = zipkintracer.NewScribeCollector(zipkinHostAndPort, time.Second, zipkintracer.ScribeBatchSize(0), zipkintracer.ScribeBatchInterval(time.Millisecond))
		if err != nil {
			log.Fatal(err)
		}
	})

	traceID += 1
	spanID += 1
	var (
		methodName   = "method"
		traceID      = traceID
		spanID       = spanID
		parentSpanID = int64(0)
		value        = message
	)

	span := makeNewSpan(methodName, traceID, spanID, parentSpanID, duration, true)
	annotate(span, time.Now(), value, nil)

	fmt.Fprint(logOut, fmt.Sprintf("[INFO] Traceid:%d spanid:%d parentid:%d\n", traceID, spanID, parentSpanID))
	return collector.Collect(span)
}

func makeNewSpan(methodName string, traceID, spanID, parentSpanID int64, duration time.Duration, debug bool) *zipkincore.Span {
	timestamp := time.Now().UnixNano() / 1e3
	return &zipkincore.Span{
		TraceID:   traceID,
		Name:      methodName,
		ID:        spanID,
		ParentID:  &parentSpanID,
		Debug:     debug,
		Duration:  thrift.Int64Ptr(duration.Nanoseconds() * 1000),
		Timestamp: &timestamp,
	}
}

func annotate(span *zipkincore.Span, timestamp time.Time, value string, host *zipkincore.Endpoint) {
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	span.Annotations = append(span.Annotations, &zipkincore.Annotation{
		Timestamp: timestamp.UnixNano() / 1e3,
		Value:     value,
		Host:      host,
	})
}

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
