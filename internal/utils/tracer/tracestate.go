// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
// Package tracer check the format of the tracestate header and extrac its sub-members;
// It expects the tracestate header is in the format of:
//  tracestate format: b3={TraceId}-{SpanId}[-{Sampled}[-{ParentId}]]

// Current version DOES NOT consider the optional case when tracestate is: b3={Sampled}
package tracer

import (
	"errors"
	"regexp"
)

var (
	// ErrInvalidFormat occurs when the format is invalid, such as if there are missing characters
	// or a field contains an unexpected character set.
	ErrInvalidFormat = errors.New("tracecontext: Invalid format")

	// ErrInvalidTraceID occurs when the encoded trace ID is invalid, i.e., all 0s
	ErrInvalidTraceID = errors.New("tracecontext: Invalid trace ID")

	// ErrInvalidSpanID occurs when the encoded span ID is invalid, i.e., all 0s
	ErrInvalidSpanID = errors.New("tracecontext: Invalid span ID")

	// ErrInvalidTraceIDSpanID occurs when the encoded trace and span IDs are invalid, i.e., all 0s
	ErrInvalidTraceIDSpanID = errors.New("tracecontext: Invalid trace ID and span ID")
)

var (
	// tracestate format: b3={TraceId}-{SpanId}[-{Sampled}[-{ParentId}]] | b3={Sampled}
	reTracestate = regexp.MustCompile(`^b3=([a-f0-9]{32})-([a-f0-9]{16})(-.*)?$`)

	// traceparent format: 00-{TraceId 16bytes}-{ParentId 8 bytes}-0{Sampled 1 Hex}
	reTraceparent = regexp.MustCompile(`^([a-f0-9]{2})-([a-f0-9]{32})-([a-f0-9]{16})-([a-f0-9]{2})(-.*)?$`)

	invalidTraceIDAllZeroes = "00000000000000000000000000000000"
	invalidSpanIDAllZeroes  = "0000000000000000"
)

// TraceState indicates information about a span and the trace of which it is part,
// so that child spans started in the same trace may propagate necessary data and share relevant behavior.
type TraceState struct {
	// TraceID is the trace ID of the whole trace, and should be constant across all spans in a given trace.
	// A `TraceID` that contains only 0 bytes should be treated as invalid.
	TraceID string

	SpanID string
}

type TraceParent struct {
	// TraceID is the trace ID of the whole trace, and should be constant across all spans in a given trace.
	// A `TraceID` that contains only 0 bytes should be treated as invalid.
	TraceID string

	ParentID string
	Sampled  bool
}

type TraceParser interface {
	ParseTraceParent(string) (TraceParent, error)
	ParseTraceState(string) (TraceState, error)
}

type TraceParserImpl struct {
}

const (
	// tracestate format: b3={TraceId}-{SpanId}[-{Sampled}[-{ParentId}]]

	// TracestateTraceIDIndex is the index of TraceID in the array of RegExp Captures for tracestate
	TracestateTraceIDIndex = 1

	// TracestateSpanIDIndex is the index of SpanID in the array of RegExp Captures for tracestate
	TracestateSpanIDIndex = 2

	// TracestateNumberOfElements is the expected number of RegExp Captures for tracestate
	TracestateNumberOfElements = 3
)

const (
	// traceparent format: 00-{TraceId 16bytes}-{ParentId 8 bytes}-0{Sampled 1 Hex}
	// TraceparentTraceIDIndex is the index of TraceID in the array of RegExp Captures for traceparent
	TraceparentTraceIDIndex = 2
	// TraceparentParentIDIndex is the index of ParentID in the array of RegExp Captures for traceparent
	TraceparentParentIDIndex = 3
	// TraceparentParentIDIndex is the index of Sampled in the array of RegExp Captures for traceparent
	TraceparentSampledIndex = 4
	// TraceparentNumberOfElements is the expected number of RegExp Captures for traceparent
	TraceparentNumberOfElements = 6 // expected 5 but works with 6
)

// ParseTraceState attempts to decode a `TraceState` from a string.
// It returns an error if the string is incorrectly formatted or otherwise invalid.
// parse b3={TraceId}-{SpanId}[-{Sampled}[-{ParentId}]]
func (t *TraceParserImpl) ParseTraceState(s string) (tp TraceState, err error) {
	err = nil

	matches := reTracestate.FindStringSubmatch(s) // matches is array of string

	// check at least these are present: full string, TraceId, SpanId
	if len(matches) < TracestateNumberOfElements {
		err = ErrInvalidFormat
		return tp, err
	}

	// extract the mandatory TraceId and SpanID
	traceID := matches[TracestateTraceIDIndex]
	if traceID == invalidTraceIDAllZeroes {
		err = ErrInvalidTraceID
	}

	spanID := matches[TracestateSpanIDIndex]
	if spanID == invalidSpanIDAllZeroes {
		if traceID == invalidTraceIDAllZeroes {
			err = ErrInvalidTraceIDSpanID
		} else {
			err = ErrInvalidSpanID
		}
	}

	tp.TraceID = traceID
	tp.SpanID = spanID

	return tp, err
}

// ParseTraceParent parses traceparent format: 00-{TraceId 16bytes}-{ParentId 8 bytes}-0{Sampled 1 Hex}
func (t *TraceParserImpl) ParseTraceParent(s string) (tp TraceParent, err error) {
	err = nil
	matches := reTraceparent.FindStringSubmatch(s)

	if len(matches) != TraceparentNumberOfElements {
		err = ErrInvalidFormat
		return tp, err
	}

	// extract the mandatory TraceId and SpanID
	traceID := matches[TraceparentTraceIDIndex]
	if traceID == invalidTraceIDAllZeroes {
		err = ErrInvalidTraceID
	}

	tp.TraceID = traceID
	tp.ParentID = matches[TraceparentParentIDIndex]
	tp.Sampled = matches[TraceparentSampledIndex] == "01"

	return tp, err
}

func NewTraceParser() TraceParser {
	return &TraceParserImpl{}
}
