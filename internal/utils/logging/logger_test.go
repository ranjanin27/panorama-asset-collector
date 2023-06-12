// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package logging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/contextutilities"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
	"google.golang.org/grpc/metadata"
)

func TestLogging(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{

		{
			name: "Logger",
			args: args{args: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logging.GetLogger().Info()
			logging.GetLogger().Infof("Unit test Info")
			logging.GetLogger().Print()
			logging.GetLogger().Printf("Unit test Trace")
			logging.GetLogger().Debug()
			logging.GetLogger().Debugf("Unit test Debug")
			logging.GetLogger().Error()
			logging.GetLogger().Errorf("Unit test Error")
			logging.GetLogger().Warn()
			logging.GetLogger().Warnf("Unit test Warn")

			logging.GetLogger().Close()
		})
	}
}

func TestGRPCLogger(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
	}{

		{
			name: "Logger",
			args: args{args: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := metadata.New(map[string]string{
				contextutilities.RequestID:    "request_id",
				contextutilities.Sampled:      "1",
				contextutilities.ParentSpanID: "parent_span_id",
				contextutilities.TraceID:      "trace_id",
				contextutilities.SpanID:       "span_id",
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			assert.NotNil(t, ctx)
			logging.GetLogger().WithContext(ctx).Info("Test Being executed")
		})
	}
}
