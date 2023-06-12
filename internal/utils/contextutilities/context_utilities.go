// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package contextutilities

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

const (
	AuthHeader = `Authorization`

	Sampled       = "x-b3-sampled"
	TraceID       = "x-b3-traceid"
	SpanID        = "x-b3-spanid"
	ParentSpanID  = "x-b3-parentspanid"
	RequestID     = "x-request-id"
	True          = "true"
	LoggerTraceID = "x-b3-traceid"
	LoggerSpanID  = "x-b3-spanid"
)

type TraceProperties struct {
	RequestID    string
	TraceID      string
	SpanID       string
	ParentSpanID string
	Sampled      bool
}

func (trace *TraceProperties) setRequestID(requestID string) *TraceProperties {
	if requestID != "" {
		trace.RequestID = requestID
	}
	return trace
}

func (trace *TraceProperties) setTraceID(traceID string) *TraceProperties {
	if traceID != "" {
		trace.TraceID = traceID
	}
	return trace
}

func (trace *TraceProperties) setSpanID(spanID string) *TraceProperties {
	if spanID != "" {
		trace.SpanID = spanID
	}
	return trace
}

func (trace *TraceProperties) setParentSpanID(parentSpanID string) *TraceProperties {
	if parentSpanID != "" {
		trace.ParentSpanID = parentSpanID
	}
	return trace
}

func (trace *TraceProperties) setSampled(sampled bool) *TraceProperties {
	trace.Sampled = sampled
	return trace
}

type supplementedMD metadata.MD

func (incoming supplementedMD) getFirst(k string) string {
	values := metadata.MD(incoming).Get(k)
	if len(values) < 1 {
		return ""
	}
	return values[0]
}

func extractGRPCTraceDetails(ctx context.Context) *TraceProperties {
	trace := new(TraceProperties)
	incomingMeta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return trace
	}
	incoming := supplementedMD(incomingMeta)
	sampled := strings.ToLower(incoming.getFirst(Sampled))
	traceID := incoming.getFirst(TraceID)
	spanID := incoming.getFirst(SpanID)
	requestID := incoming.getFirst(RequestID)
	parentSpanID := incoming.getFirst(ParentSpanID)
	isSampled := sampled == "1" || sampled == True
	trace.setTraceID(traceID).setSpanID(spanID).setRequestID(requestID).setParentSpanID(parentSpanID).setSampled(isSampled)
	return trace
}

func extractHTTPTraceDetails(ctx *gin.Context) *TraceProperties {
	req := ctx.Request
	trace := new(TraceProperties)
	sampled := strings.ToLower(req.Header.Get(Sampled))
	traceID := req.Header.Get(TraceID)
	spanID := req.Header.Get(SpanID)
	requestID := req.Header.Get(RequestID)
	parentSpanID := req.Header.Get(ParentSpanID)
	isSampled := sampled == "1" || sampled == True
	trace.setParentSpanID(parentSpanID).setRequestID(requestID).setSampled(isSampled).setTraceID(traceID).setSpanID(spanID)
	return trace
}

func ExtractTraceDetails(ctx context.Context) *TraceProperties {
	var trace *TraceProperties
	if c, ok := ctx.(*gin.Context); ok {
		trace = extractHTTPTraceDetails(c)
	} else {
		trace = extractGRPCTraceDetails(ctx)
	}
	return trace
}

func (trace *TraceProperties) AppendToOutgoingContext(ctx context.Context) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = trace.addToHTTPContext(c)
	} else {
		ctx = trace.appendToOutgoingGRPCContext(ctx)
	}
	return ctx
}

func (trace *TraceProperties) appendToOutgoingGRPCContext(ctx context.Context) context.Context {
	if trace.Sampled {
		ctx = metadata.AppendToOutgoingContext(ctx, Sampled, "1")
	} else {
		ctx = metadata.AppendToOutgoingContext(ctx, Sampled, "0")
	}
	ctx = metadata.AppendToOutgoingContext(ctx, TraceID, trace.TraceID)
	ctx = metadata.AppendToOutgoingContext(ctx, SpanID, trace.SpanID)
	ctx = metadata.AppendToOutgoingContext(ctx, ParentSpanID, trace.ParentSpanID)
	ctx = metadata.AppendToOutgoingContext(ctx, RequestID, trace.RequestID)
	return ctx
}

func (trace *TraceProperties) addToHTTPContext(c *gin.Context) *gin.Context {
	if c == nil {
		return nil
	}
	req := c.Request
	if req == nil {
		req = &http.Request{}
		req.Header = http.Header{}
	}
	if trace.Sampled {
		req.Header.Add(Sampled, "1")
	} else {
		req.Header.Add(Sampled, "0")
	}
	req.Header.Add(TraceID, trace.TraceID)
	req.Header.Add(SpanID, trace.SpanID)
	req.Header.Add(ParentSpanID, trace.ParentSpanID)
	req.Header.Add(RequestID, trace.RequestID)

	c.Request = req
	return c
}
