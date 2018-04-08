package main

import (
	"time"
	"github.com/openzipkin/zipkin-go-opentracing"
	"math/rand"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/apache/thrift/lib/go/thrift"
)

var collector, _ = zipkintracer.NewScribeCollector("localhost:9410", time.Second, zipkintracer.ScribeBatchSize(0), zipkintracer.ScribeBatchInterval(time.Millisecond))

var traceID = int64(rangeIn(100000, 999999))
var spanID = int64(rangeIn(100000, 999999))

func trace(message string, duration time.Duration) error {
	traceID += 1
	var (
		methodName   = "method"
		traceID      = traceID
		spanID       = spanID
		parentSpanID = int64(0)
		value        = message
	)

	span := makeNewSpan(methodName, traceID, spanID, parentSpanID, duration, true)
	annotate(span, time.Now(), value, nil)

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
