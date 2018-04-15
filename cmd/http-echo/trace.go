package main

import (
	"fmt"
	"github.com/openzipkin/zipkin-go-opentracing"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/apache/thrift/lib/go/thrift"
	"math/rand"
	"net/http"
	"encoding/base64"
	"github.com/cenkalti/backoff"
)

var collector zipkintracer.Collector

var traceID = int64(rangeIn(100000, 999999))
var spanID = int64(rangeIn(100000, 999999))
var once sync.Once
var zipkinHostAndPort = fmt.Sprintf("%s:%s",
	valueOrDefault(os.Getenv("ZIPKIN_HOST"), "localhost"),
	valueOrDefault(os.Getenv("ZIPKIN_PORT"), "9410"))

func trace(request *http.Request, logOut io.Writer, message string, duration time.Duration) error {
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

/*
  final String trace64 = httpServletRequest.getHeader("l5d-ctx-trace");

            if(StringUtils.isNotEmpty(trace64)) {
                final byte[] traceBytes = Base64.getDecoder().decode(trace64);
                val trace = TraceId.deserialize(traceBytes).get();

                val spanId = trace.spanId().toString();
                val traceId = trace.traceId().toString();
                val parentId = trace.parentId().toString();

                String sampled = "0";
                if(!trace.sampled().isEmpty() && trace.sampled().get() instanceof Boolean && (Boolean)trace.sampled().get()) {
                    sampled = "1";
                }

                chain.doFilter(new SpanHttpServletRequest(httpServletRequest, spanId, traceId, parentId, sampled), response);
                return;
            }
 */
func getParentSpanId(request *http.Request) (error) {
	traceHeaders, ok := request.Header["l5d-ctx-trace"]
	if !ok {
		return nil
	}

	_, err := base64.StdEncoding.DecodeString(traceHeaders[0])
	if err != nil {

		return err
	}

	return nil
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

func retry(op func() error, timeout time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second * 5
	bo.MaxElapsedTime = timeout
	return backoff.Retry(op, bo)
}
