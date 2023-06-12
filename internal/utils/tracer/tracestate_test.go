// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceParserImpl_ParseTraceState(t *testing.T) {
	// Test with valid tracestate header
	tracestateHeader := "b3=a8a090a7ba3e6499345a60106220c295-854bfaaa86165a5a"

	TraceParser := NewTraceParser()
	traceState, err := TraceParser.ParseTraceState(tracestateHeader)

	assert.Nil(t, err, ErrInvalidFormat)
	assert.Equal(t, "a8a090a7ba3e6499345a60106220c295", traceState.TraceID)
	assert.Equal(t, "854bfaaa86165a5a", traceState.SpanID)

	// Test with valid tracestateheader including optional headers (sampled and parentID)
	tracestateHeader = "b3=a8a090a7ba3e6499345a60106220c295-854bfaaa86165a5a-0-9a388e1a73d0701d"
	TraceParser = NewTraceParser()
	traceState, err = TraceParser.ParseTraceState(tracestateHeader)

	assert.Nil(t, err, ErrInvalidFormat)
	assert.Equal(t, "a8a090a7ba3e6499345a60106220c295", traceState.TraceID)
	assert.Equal(t, "854bfaaa86165a5a", traceState.SpanID)

	// Test with invalid tracestate header with all fields
	tracestateHeader = "b3=00000000000000000000000000000000-0000000000000000-0-0000000000000000"
	traceState, err = TraceParser.ParseTraceState(tracestateHeader)
	assert.Equal(t, err, ErrInvalidTraceIDSpanID)
	assert.Equal(t, "00000000000000000000000000000000", traceState.TraceID)
	assert.Equal(t, "0000000000000000", traceState.SpanID)

	// Test with tracestate header with invalid spanID
	tracestateHeader = "b3=a8a090a7ba3e6499345a60106220c295-0000000000000000-0-9a388e1a73d0701d"
	traceState, err = TraceParser.ParseTraceState(tracestateHeader)
	assert.Equal(t, err, ErrInvalidSpanID)
	assert.Equal(t, "a8a090a7ba3e6499345a60106220c295", traceState.TraceID)
	assert.Equal(t, "0000000000000000", traceState.SpanID)

	// Test with invalid tracestate header
	tracestateHeader = "b3=00000000000000000000000000000000"
	traceState, err = TraceParser.ParseTraceState(tracestateHeader)
	assert.Equal(t, err, ErrInvalidFormat)
}

func TestTraceParserImpl_ParseTraceParent(t *testing.T) {
	// traceparent format: 00-{TraceId 16bytes}-{ParentId 8 bytes}-0{Sampled 1 Hex}
	// Test with valid traceparent header
	header := "00-a8a090a7ba3e6499345a60106220c295-854bfaaa86165a5a-01"
	parser := NewTraceParser()
	tp, err := parser.ParseTraceParent(header)

	assert.Nil(t, err, ErrInvalidFormat)
	assert.Equal(t, "a8a090a7ba3e6499345a60106220c295", tp.TraceID)
	assert.Equal(t, "854bfaaa86165a5a", tp.ParentID)
	assert.True(t, tp.Sampled)

	// Test with valid traceparent header with sampled == false
	header = "00-a8a090a7ba3e6499345a60106220c295-854bfaaa86165a5a-00"
	tp, err = parser.ParseTraceParent(header)
	assert.Nil(t, err, ErrInvalidFormat)
	assert.False(t, tp.Sampled)

	// Test with invalid traceparent header
	header = "00-00000000000000000000000000000000-0000000000000000-00"
	tp, err = parser.ParseTraceParent(header)
	assert.Equal(t, err, ErrInvalidTraceID)

	// Test with invalid traceparent header
	header = "00000000000000000000000000000000"
	tp, err = parser.ParseTraceParent(header)
	assert.Equal(t, err, ErrInvalidFormat)
}
