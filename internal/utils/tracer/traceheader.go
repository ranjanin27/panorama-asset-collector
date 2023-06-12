// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package tracer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/constants"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
	"google.golang.org/grpc/metadata"
)

const True = "true"

const (
	prefix           = "aaaaaaaa"
	traceidSuffixLen = 24
	spanidSuffixLen  = 8
)

const (
	GRPCSampled      = "x-b3-sampled"
	GRPCTraceID      = "x-b3-traceid"
	GRPCSpanID       = "x-b3-spanid"
	GRPCParentSpanID = "x-b3-parentspanid"
	GRPCRequestID    = "x-request-id"

	HTTPSampled      = "X-B3-Sampled"
	HTTPTraceID      = "X-B3-TraceId"
	HTTPSpanID       = "X-B3-SpanId"
	HTTPParentSpanID = "X-B3-ParentSpanId"
	HTTPRequestID    = "X-Request-Id"

	// Logging fields are lower case so are identical to GRPC fields
	LoggerTraceID = GRPCTraceID
	LoggerSpanID  = GRPCSpanID

	CETraceParent = "ce_traceparent"
	CETraceState  = "ce_tracestate"
	TraceHeader   = "b3="
)

type TraceHeaders struct {
	RequestID    string
	TraceID      string
	SpanID       string
	ParentSpanID string
	Sampled      bool
}

type GrpcTracePropagator interface {
	AppendToOutgoingGrpcContext(context.Context) context.Context
}

type HTTPTracePropagator interface {
	AddToHTTPRequest(*http.Request) *http.Request
	GetTraceID() string
	GetSpanID() string
}

type LoggerTracePropagator interface {
	AddToLogger(logger logging.Logger) logging.Logger
}

type TracePropagator interface {
	GrpcTracePropagator
	HTTPTracePropagator
	LoggerTracePropagator
}

func (hdrs *TraceHeaders) AddToLogger(logger logging.Logger) logging.Logger {
	return logger.WithFields(LoggerTraceID, hdrs.TraceID, LoggerSpanID, hdrs.SpanID)
}

func (hdrs *TraceHeaders) AppendToOutgoingGrpcContext(ctx context.Context) context.Context {
	if hdrs.Sampled {
		ctx = metadata.AppendToOutgoingContext(ctx, GRPCSampled, "1")
	} else {
		ctx = metadata.AppendToOutgoingContext(ctx, GRPCSampled, "0")
	}
	ctx = metadata.AppendToOutgoingContext(ctx, GRPCTraceID, hdrs.TraceID)
	ctx = metadata.AppendToOutgoingContext(ctx, GRPCSpanID, hdrs.SpanID)
	ctx = metadata.AppendToOutgoingContext(ctx, GRPCParentSpanID, hdrs.ParentSpanID)
	ctx = metadata.AppendToOutgoingContext(ctx, GRPCRequestID, hdrs.RequestID)
	return ctx
}

func (hdrs *TraceHeaders) AddToHTTPRequest(req *http.Request) *http.Request {
	if hdrs.Sampled {
		req.Header.Add(HTTPSampled, "1")
	} else {
		req.Header.Add(HTTPSampled, "0")
	}
	req.Header.Add(HTTPTraceID, hdrs.TraceID)
	req.Header.Add(HTTPSpanID, hdrs.SpanID)
	req.Header.Add(HTTPParentSpanID, hdrs.ParentSpanID)
	req.Header.Add(HTTPRequestID, hdrs.RequestID)
	return req
}

func (hdrs *TraceHeaders) GetTraceID() string {
	return hdrs.TraceID
}

func (hdrs *TraceHeaders) GetSpanID() string {
	return hdrs.SpanID
}

func (hdrs *TraceHeaders) TraceParent() string {
	sd := "0"
	if hdrs.Sampled {
		sd = "1"
	}
	tp := fmt.Sprintf(
		"00-%s-%s-0%s",
		hdrs.TraceID,
		hdrs.ParentSpanID,
		sd,
	)
	return tp
}

func (hdrs *TraceHeaders) TraceState() string {
	ts := TraceHeader + "0"
	if hdrs.Sampled {
		ts = fmt.Sprintf(
			TraceHeader+"%s-%s-1-%s",
			hdrs.TraceID,
			hdrs.SpanID,
			hdrs.ParentSpanID,
		)
	}
	return ts
}

func ExtractFromHTTPRequest(req *http.Request, logger logging.Logger) TracePropagator {
	hdrSampled := strings.ToLower(req.Header.Get(HTTPSampled))
	hdrs := &TraceHeaders{
		TraceID:      req.Header.Get(HTTPTraceID),
		SpanID:       req.Header.Get(HTTPSpanID),
		ParentSpanID: req.Header.Get(HTTPParentSpanID),
		RequestID:    req.Header.Get(HTTPRequestID),
		Sampled:      hdrSampled == "1" || hdrSampled == True,
	}
	if ok := doHeadersExist(hdrs); !ok {
		addMissingHeaders(hdrs)
		logger.WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received REST request with one or more empty required tracer headers. Adding TraceId and SpanId now")
	} else {
		logger.WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).Info("Received REST request with TraceID: ")
	}

	return hdrs
}

func ExtractFromIncomingGrpcContext(ctx context.Context, logger logging.Logger) TracePropagator {
	incomingMeta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	incoming := supplementedMD(incomingMeta)
	hdrSampled := strings.ToLower(incoming.getFirst(GRPCSampled))
	hdrs := &TraceHeaders{
		TraceID:      incoming.getFirst(GRPCTraceID),
		SpanID:       incoming.getFirst(GRPCSpanID),
		ParentSpanID: incoming.getFirst(GRPCParentSpanID),
		RequestID:    incoming.getFirst(GRPCRequestID),
		Sampled:      hdrSampled == "1" || hdrSampled == True,
	}
	if ok := doHeadersExist(hdrs); !ok {
		addMissingHeaders(hdrs)
		logger.
			WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received gRPC request with one or more empty required tracer headers. Adding TraceId and SpanId now")
	} else {
		logger.WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).Info("Received gRPC request with TraceID: ")
	}
	return hdrs
}

func ExtractHeadersFromIncomingGrpcContext(ctx context.Context, logger logging.Logger) TraceHeaders {
	hdrs := &TraceHeaders{}
	incomingMeta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return *hdrs
	}
	incoming := supplementedMD(incomingMeta)
	hdrSampled := strings.ToLower(incoming.getFirst(GRPCSampled))
	hdrs = &TraceHeaders{
		TraceID:      incoming.getFirst(GRPCTraceID),
		SpanID:       incoming.getFirst(GRPCSpanID),
		ParentSpanID: incoming.getFirst(GRPCParentSpanID),
		RequestID:    incoming.getFirst(GRPCRequestID),
		Sampled:      hdrSampled == "1" || hdrSampled == True,
	}
	if ok := doHeadersExist(hdrs); !ok {
		addMissingHeaders(hdrs)
		logger.
			WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received gRPC request with one or more empty required tracer headers. Adding TraceId and SpanId now")
	} else {
		logger.WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).Info("Received gRPC request with TraceID: ")
	}
	return *hdrs
}

func ExtractFromIncomingKafkaEvent(eventData map[string]string, logger logging.Logger) *TraceHeaders {
	hdrs := &TraceHeaders{}
	traceParser := NewTraceParser()

	// Check TraceParent first: https://w3c.github.io/trace-context/#tracestate-header
	// If TraceParent cannot be parsed, spec says to not try to parse the TraceState
	tpString, ok := eventData[CETraceParent]
	if !ok {
		addMissingHeaders(hdrs)
		logger.
			WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received Kafka Event with invalid traceparent")
		return hdrs
	}

	tp, err := traceParser.ParseTraceParent(tpString)
	if err != nil {
		addMissingHeaders(hdrs)
		logger.
			WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received kafka Event with invalid traceparent")
		return hdrs
	}
	hdrs.TraceID = tp.TraceID
	hdrs.ParentSpanID = tp.ParentID
	hdrs.Sampled = tp.Sampled

	// Extract TraceID & SpanID from TraceState
	tsString, ok := eventData[CETraceState]
	if !ok {
		addMissingHeaders(hdrs)
		logger.
			WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received kafka Event with invalid tracestate")
		return hdrs
	}

	ts, err := traceParser.ParseTraceState(tsString)
	if err != nil {
		addMissingHeaders(hdrs)
		logger.
			WithFields(LoggerTraceID, hdrs.TraceID).
			WithFields(LoggerSpanID, hdrs.SpanID).
			Warn("received kafka Event with invalid tracestate")
		return hdrs
	}
	hdrs.SpanID = ts.SpanID

	// Note: RequestID header is not applicable for kafka events
	return hdrs
}

// Private types and methods

type supplementedMD metadata.MD

func (incoming supplementedMD) getFirst(k string) string {
	values := metadata.MD(incoming).Get(k)
	if len(values) < 1 {
		return ""
	}
	return values[0]
}

func doHeadersExist(
	hdrs *TraceHeaders,
) bool {
	return hdrs.TraceID != "" && hdrs.SpanID != ""
}

func addMissingHeaders(hdrs *TraceHeaders) {
	if hdrs.TraceID == "" {
		hdrs.TraceID = generateHeader(traceidSuffixLen)
	}
	if hdrs.SpanID == "" {
		hdrs.SpanID = generateHeader(spanidSuffixLen)
	}
	if hdrs.ParentSpanID == "" {
		hdrs.ParentSpanID = hdrs.SpanID
	}
}

func generateHeader(suffixLen int) string {
	return prefix + hexString(suffixLen)
}

func hexString(digits int) string {
	// How many bytes of data to randomly generate
	byteLen := digits >> 1
	if digits&1 == 1 {
		byteLen++
	}
	// Allocate and fill scratch space
	buf := make([]byte, byteLen)
	rand.Read(buf) //nolint:gosec // don't need this to be crypto-secure
	// Convert to hex string and return
	hexStr := hex.EncodeToString(buf)[0:digits]
	return strings.ToLower(hexStr)
}

// MakeCETraceParent constructs the traceparent string needed for tracing with Kafka
//   Returns an error as part of the prototype in case needed later
//   (e.g. if validation is performed on the headers)

func MakeCETraceParent(headers *TraceHeaders) string {
	sd := "0"
	if headers.Sampled {
		sd = "1"
	}
	output := fmt.Sprintf(
		"00-%s-%s-0%s",
		headers.TraceID,
		headers.ParentSpanID,
		sd,
	)
	return output
}

func MakeCETraceState(headers *TraceHeaders) string {
	output := constants.TraceStatePrefix
	// If Sampled header is not present, sampling decision has been deferred.
	//   Should be treated as true
	//   See: https://pages.github.hpe.com/cloud/storage-design/docs/observability.html#tracing-headers
	// "false" option is to cover legacy implementations - will be propagated as 0 or 1
	if headers.Sampled {
		output = fmt.Sprintf(
			"b3=%s-%s-1-%s",
			headers.TraceID,
			headers.SpanID,
			headers.ParentSpanID,
		)
	}
	return output
}
