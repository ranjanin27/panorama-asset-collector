// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package tracer

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

const (
	TestDataDir  = "testdata"
	GrpcTestFile = "grpc_test_header.json"
	HTTPTestFile = "http_test_header.json"
)

const (
	expectedParent = "00-trace_id-parent_span_id-01"
	expectedState  = "b3=trace_id-span_id-1-parent_span_id"
)

var (
	gth = &TraceHeaders{
		TraceID:      "trace_id",
		SpanID:       "span_id",
		ParentSpanID: "parent_span_id",
		RequestID:    "request_id",
		Sampled:      true,
	}
)

var (
	logger = logging.GetLogger()
)

func AssertTraceHeaderEquality(t *testing.T, expected, actual *TraceHeaders, msg string) {
	assertify := assert.New(t)
	if assertify.NotNilf(actual, "%s: actual is nil", msg) {
		assertify.Equalf(expected.RequestID, actual.RequestID, "%s: RequestID does not match", msg)
		assertify.Equalf(expected.TraceID, actual.TraceID, "%s: TraceID does not match", msg)
		assertify.Equalf(expected.SpanID, actual.SpanID, "%s: SpanID does not match", msg)
		assertify.Equalf(expected.ParentSpanID, actual.ParentSpanID, "%s: ParentSpanID does not match", msg)
		assertify.Equalf(expected.Sampled, actual.Sampled, "%s: Sampled does not match", msg)
	}
}

func TestTraceHeaders(t *testing.T) {
	th := gth

	expectedTP := "00-trace_id-parent_span_id-01"
	tP := th.TraceParent()
	assert.Equal(t, expectedTP, tP)

	expectedTS := "b3=trace_id-span_id-1-parent_span_id"
	tS := th.TraceState()
	assert.Equal(t, expectedTS, tS)

	th.Sampled = false

	expectedTP = "00-trace_id-parent_span_id-00"
	tP = th.TraceParent()
	assert.Equal(t, expectedTP, tP)

	expectedTS = "b3=0"
	tS = th.TraceState()
	assert.Equal(t, expectedTS, tS)

	th.Sampled = true
}

func TestTraceHeader_HTTPExtraction(t *testing.T) {
	var cases []struct {
		Name    string
		Headers map[string]string
		Want    TraceHeaders
	}
	NewTestDataInterface(TestDataDir, t).MustGetJSON(HTTPTestFile, &cases)
	for _, c := range cases {
		c := c
		t.Run(c.Name, func(tt *testing.T) {
			assertify := assert.New(tt)

			// Arrange
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, "GET", "", http.NoBody)
			assertify.NoError(err)
			assertify.NotNil(req, "expected non-nil return from NewRequest")

			for k, v := range c.Headers {
				req.Header.Add(k, v)
			}

			// Act
			extractedPropagator := ExtractFromHTTPRequest(req, logger)

			// Assert
			assertify.NotNil(extractedPropagator, "expected non-nil return from ExtractFromHttpRequest")

			extracted, ok := extractedPropagator.(*TraceHeaders)
			if assertify.True(ok, "should return TraceHeaders struct") {
				AssertTraceHeaderEquality(tt, &c.Want, extracted, "expected correct return from ExtractFromHttpRequest")
			}
		})
	}
}

func Test_ExtractFromHTTPRequest_MissingTraceHeaders(t *testing.T) {
	assertify := assert.New(t)
	requireify := require.New(t)

	cases := []struct {
		name    string
		traceID string
		spanID  string
	}{
		{
			name:    "both empty",
			traceID: "",
			spanID:  "",
		},
		{
			name:    "traceID empty",
			traceID: "",
			spanID:  "span-id",
		},
		{
			name:    "spanID empty",
			traceID: "trace-id",
			spanID:  "",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(tt *testing.T) {
			// Arrange
			hdrs := TraceHeaders{
				TraceID: c.traceID,
				SpanID:  c.spanID,
			}
			req, err := http.NewRequest("GET", "", http.NoBody)
			requireify.NoError(err)
			requireify.NotNil(req, "expected non-nil return from NewRequest")

			req = hdrs.AddToHTTPRequest(req)

			// Act
			extractedPropagator := ExtractFromHTTPRequest(req, logger)

			// Assert
			assertify.NotNil(extractedPropagator, "expected non-nil return from ExtractFromHttpRequest")

			extracted, ok := extractedPropagator.(*TraceHeaders)
			requireify.True(ok, "type deduction failed on extractedPropagator")

			if c.traceID == "" {
				assertify.Len(extracted.TraceID, 32)
				assertify.Equal("aaaaaaaa", extracted.TraceID[0:8])
			} else {
				assertify.Equal(c.traceID, extracted.TraceID)
			}

			if c.spanID == "" {
				assertify.Len(extracted.SpanID, 16)
				assertify.Equal("aaaaaaaa", extracted.SpanID[0:8])
			} else {
				assertify.Equal(c.spanID, extracted.SpanID)
			}
		})
	}
}

func TestExtractFromIncomingGrpc(t *testing.T) {
	md := metadata.New(nil)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	th := ExtractFromIncomingGrpcContext(ctx, logging.GetLogger())
	assert.NotNil(t, th)

	md = metadata.New(map[string]string{
		GRPCRequestID:    "request_id",
		GRPCSampled:      "1",
		GRPCParentSpanID: "parent_span_id",
		GRPCTraceID:      "trace_id",
		GRPCSpanID:       "span_id",
	})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	th = ExtractFromIncomingGrpcContext(ctx, logging.GetLogger())
	assert.Equal(t, gth, th)

	hdrs := ExtractHeadersFromIncomingGrpcContext(ctx, logging.GetLogger())
	hdrs.GetTraceID()
	hdrs.GetSpanID()
	assert.NotNil(t, hdrs)
}

func TestTraceHeader_GRPCExtraction(t *testing.T) {
	var cases []struct {
		Name    string
		Headers map[string]string
		Want    TraceHeaders
	}
	NewTestDataInterface(TestDataDir, t).MustGetJSON(GrpcTestFile, &cases)
	for _, c := range cases {
		c := c
		t.Run(c.Name, func(tt *testing.T) {
			assertify := assert.New(tt)

			// Arrange
			incomingMeta := metadata.New(c.Headers)
			assertify.NotNil(incomingMeta, "expected non-nil return from metadata.New")

			incomingCtx := metadata.NewIncomingContext(context.Background(), incomingMeta)
			assertify.NotNil(incomingCtx, "expected non-nil return from NewIncomingContext")

			// Act
			extractedPropagator := ExtractFromIncomingGrpcContext(incomingCtx, logger)

			// Assert
			assertify.NotNil(extractedPropagator, "expected non-nil return from ExtractFromIncomingGrpcContext")

			extracted, ok := extractedPropagator.(*TraceHeaders)
			if assertify.True(ok, "should return TraceHeaders struct") {
				AssertTraceHeaderEquality(tt, &c.Want, extracted, "expected correct return from ExtractFromIncomingGrpcContext")
			}
		})
	}
}

func Test_ExtractFromIncomingGrpcContext_MissingTraceHeaders(t *testing.T) {
	assertify := assert.New(t)
	requireify := require.New(t)

	cases := []struct {
		name    string
		traceID string
		spanID  string
	}{
		{
			name:    "both empty",
			traceID: "",
			spanID:  "",
		},
		{
			name:    "traceID empty",
			traceID: "",
			spanID:  "span-id",
		},
		{
			name:    "spanID empty",
			traceID: "trace-id",
			spanID:  "",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(tt *testing.T) {
			// Arrange
			hdrs := TraceHeaders{
				TraceID: c.traceID,
				SpanID:  c.spanID,
			}
			outgoingCtx := hdrs.AppendToOutgoingGrpcContext(context.Background())
			incomingCtx := MakeIncomingContext(outgoingCtx)

			// Act
			extractedPropagator := ExtractFromIncomingGrpcContext(incomingCtx, logger)

			// Assert
			assertify.NotNil(extractedPropagator, "expected non-nil return from ExtractFromHttpRequest")

			extracted, ok := extractedPropagator.(*TraceHeaders)
			requireify.True(ok, "type deduction failed on extractedPropagator")

			if c.traceID == "" {
				assertify.Len(extracted.TraceID, 32)
				assertify.Equal("aaaaaaaa", extracted.TraceID[0:8])
			} else {
				assertify.Equal(c.traceID, extracted.TraceID)
			}

			if c.spanID == "" {
				assertify.Len(extracted.SpanID, 16)
				assertify.Equal("aaaaaaaa", extracted.SpanID[0:8])
			} else {
				assertify.Equal(c.spanID, extracted.SpanID)
			}
		})
	}
}

func TestExtractFromIncomingKafka(t *testing.T) {
	th := ExtractFromIncomingKafkaEvent(nil, logging.GetLogger())
	assert.NotNil(t, th)

	expectedTH := &TraceHeaders{
		TraceID:      "80f198ee56343ba864fe8b2a57d3eff7",
		ParentSpanID: "05e3ac9a4f6e3b90",
		SpanID:       "e457b5a2e4d86bd1",
		Sampled:      true,
	}
	eventData := map[string]string{
		CETraceParent: "00-80f198ee56343ba864fe8b2a57d3eff7-05e3ac9a4f6e3b90-01",
		CETraceState:  "b3=80f198ee56343ba864fe8b2a57d3eff7-e457b5a2e4d86bd1-1-05e3ac9a4f6e3b90",
	}
	th = ExtractFromIncomingKafkaEvent(eventData, logging.GetLogger())
	assert.Equal(t, expectedTH, th)

	eventData = map[string]string{
		CETraceParent:        "00-80f198ee56343ba864fe8b2a57d3eff7-05e3ac9a4f6e3b90-01",
		"wrong-state-header": "b3=80f198ee56343ba864fe8b2a57d3eff7-e457b5a2e4d86bd1-1-05e3ac9a4f6e3b90",
	}
	th = ExtractFromIncomingKafkaEvent(eventData, logging.GetLogger())
	assert.NotNil(t, th)
}

func MakeIncomingContext(outgoingCtx context.Context) context.Context {
	outgoingMeta, ok := metadata.FromOutgoingContext(outgoingCtx)
	if !ok {
		return nil
	}
	return metadata.NewIncomingContext(context.Background(), outgoingMeta)
}

// Test_MakeCETraceParent checks the CE TraceParent is constructed correctly
func Test_MakeCETraceParent(t *testing.T) {
	traceParent := MakeCETraceParent(gth)
	assert.Equal(t, expectedParent, traceParent)
}

// Test_MakeCETraceParent checks the CE TraceParent is constructed correctly
func Test_MakeCETraceState(t *testing.T) {
	traceState := MakeCETraceState(gth)
	assert.Equal(t, expectedState, traceState)
}
