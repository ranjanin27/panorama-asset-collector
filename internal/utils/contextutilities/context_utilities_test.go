// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package contextutilities

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

const (
	TestDataDir  = "testdata"
	GrpcTestFile = "grpc_test_header.json"
	HTTPTestFile = "http_test_header.json"
)

var (
	gth = &TraceProperties{
		RequestID:    "request_id",
		TraceID:      "trace_id",
		SpanID:       "span_id",
		ParentSpanID: "parent_span_id",
		Sampled:      true,
	}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type TestData struct {
	t       *testing.T
	dataDir string
}

func TestExtractFromIncomingGrpc(t *testing.T) {
	md := metadata.New(nil)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	th := extractGRPCTraceDetails(ctx)
	assert.NotNil(t, th)

	md = metadata.New(map[string]string{
		RequestID:    "request_id",
		Sampled:      "1",
		ParentSpanID: "parent_span_id",
		TraceID:      "trace_id",
		SpanID:       "span_id",
	})
	ctx = metadata.NewIncomingContext(context.Background(), md)
	th = extractGRPCTraceDetails(ctx)
	assert.Equal(t, gth, th)
}

func TestTraceHeaderExtraction(t *testing.T) {
	var cases []struct {
		Name    string
		Headers map[string]string
		Want    TraceProperties
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
			extractedTrace := ExtractTraceDetails(incomingCtx)

			// Assert
			assertify.NotNil(extractedTrace, "expected non-nil return from ExtractFromIncomingGrpcContext")
			AssertTraceHeaderEquality(tt, &c.Want, extractedTrace, "expected correct return from ExtractFromIncomingGrpcContext")
		})
	}

	NewTestDataInterface(TestDataDir, t).MustGetJSON(HTTPTestFile, &cases)
	for _, c := range cases {
		c := c
		t.Run(c.Name, func(tt *testing.T) {
			assertify := assert.New(tt)

			// Arrange
			req, err := http.NewRequest(http.MethodGet, "", http.NoBody)
			assertify.NoError(err)
			assertify.NotNil(req, "expected non-nil return from NewRequest")
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			for k, v := range c.Headers {
				req.Header.Add(k, v)
			}
			ctx.Request = req

			// Act
			extractedTrace := ExtractTraceDetails(ctx)
			// Assert
			assertify.NotNil(&extractedTrace, "expected non-nil return from ExtractFromHttpRequest")
			AssertTraceHeaderEquality(tt, &c.Want, extractedTrace, "expected correct return from ExtractFromHttpRequest")
		})
	}
}

func Test_AppendToOutgoingContext(t *testing.T) {
	var cases []struct {
		Name    string
		Headers map[string]string
		Want    TraceProperties
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
			extractedTrace := ExtractTraceDetails(incomingCtx)

			// Assert
			assertify.NotNil(&extractedTrace, "expected non-nil return from ExtractTraceDetails")
			AssertTraceHeaderEquality(tt, &c.Want, extractedTrace, "expected correct return from ExtractTraceDetails")

			outgoingCtx := extractedTrace.AppendToOutgoingContext(context.Background())
			outgoingMeta, ok := metadata.FromOutgoingContext(outgoingCtx)
			assertify.True(ok, "Extraction from outgoing context successful")
			outgoing := supplementedMD(outgoingMeta)
			newExtractedTrace := new(TraceProperties)
			sampled := strings.ToLower(outgoing.getFirst(Sampled))
			traceID := outgoing.getFirst(TraceID)
			spanID := outgoing.getFirst(SpanID)
			requestID := outgoing.getFirst(RequestID)
			parentSpanID := outgoing.getFirst(ParentSpanID)
			isSampled := sampled == "1" || sampled == True
			newExtractedTrace.setTraceID(traceID).setSpanID(spanID).setRequestID(requestID).setParentSpanID(parentSpanID).setSampled(isSampled)

			assertify.NotNil(newExtractedTrace, "expected non-nil return from ExtractTraceDetails from outgoing context")
			AssertTraceHeaderEquality(tt, &c.Want, newExtractedTrace, "expected correct return from AppendToOutgoingContext")
		})
	}

	NewTestDataInterface(TestDataDir, t).MustGetJSON(HTTPTestFile, &cases)
	for _, c := range cases {
		c := c
		t.Run(c.Name, func(tt *testing.T) {
			assertify := assert.New(tt)

			// Arrange
			req, err := http.NewRequest(http.MethodGet, "", http.NoBody)
			assertify.NoError(err)
			assertify.NotNil(req, "expected non-nil return from NewRequest")
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

			for k, v := range c.Headers {
				req.Header.Add(k, v)
			}
			ctx.Request = req
			extractedTrace := ExtractTraceDetails(ctx)
			outgoingCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
			outgoingCtx, _ = extractedTrace.AppendToOutgoingContext(outgoingCtx).(*gin.Context)
			outgoingTrace := ExtractTraceDetails(outgoingCtx)
			assertify.NotNil(&outgoingTrace, "expected non-nil return from ExtractTraceDetails from outgoing context")
			AssertTraceHeaderEquality(tt, &c.Want, outgoingTrace, "expected correct return from AppendToOutgoingContext")
		})
	}
}

func NewTestDataInterface(dir string, t *testing.T) *TestData {
	return &TestData{
		t:       t,
		dataDir: dir,
	}
}

func AssertTraceHeaderEquality(t *testing.T, expected, actual *TraceProperties, msg string) {
	assertify := assert.New(t)
	if assertify.NotNilf(actual, "%s: actual is nil", msg) {
		assertify.Equalf(expected.RequestID, actual.RequestID, "%s: RequestID does not match", msg)
		assertify.Equalf(expected.TraceID, actual.TraceID, "%s: TraceID does not match", msg)
		assertify.Equalf(expected.SpanID, actual.SpanID, "%s: SpanID does not match", msg)
		assertify.Equalf(expected.ParentSpanID, actual.ParentSpanID, "%s: ParentSpanID does not match", msg)
		assertify.Equalf(expected.Sampled, actual.Sampled, "%s: Sampled does not match", msg)
	}
}

func (td *TestData) MustGetJSON(filename string, contents interface{}) {
	byteSequence := td.MustGetContents(filename)
	err := json.Unmarshal(byteSequence, contents)
	if err != nil {
		td.t.Fatal(err)
	}
}

func (td *TestData) MustGetContents(filename string) []byte {
	path := filepath.Join(td.dataDir, filename)
	file, err := os.Open(path)
	if err != nil {
		td.t.Fatal(err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			td.t.Fatal(closeErr)
		}
	}()
	byteSequence, err := io.ReadAll(file)
	if err != nil {
		td.t.Fatal(err)
	}
	return byteSequence
}
