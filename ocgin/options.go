package ocgin

import (
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

type TraceOption func(o *TraceOptions)

type TraceOptions struct {
	// DefaultAttributes will be set to each span as default.
	DefaultAttributes []trace.Attribute

	// IsPublicEndpoint should be set to true for publicly accessible HTTP(S)
	// servers. If true, any trace metadata set on the incoming request will
	// be added as a linked trace instead of being added as a parent of the
	// current trace.
	IsPublicEndpoint bool

	// Sample is applied to the span started by this Handler around each
	// request. default is 0.0001
	Sample func(c *gin.Context) trace.Sampler
}

var _defaultOptions = TraceOptions{
	DefaultAttributes: []trace.Attribute{},
	IsPublicEndpoint:  false,
	Sample: func(c *gin.Context) trace.Sampler {
		return trace.ProbabilitySampler(0.001)
	},
}

// WithDefaultAttributes will be set to each span as default.
func WithDefaultAttributes(attrs ...trace.Attribute) TraceOption {
	return func(o *TraceOptions) {
		o.DefaultAttributes = attrs
	}
}

func WithPublicEndpoint(isPublic bool) TraceOption {
	return func(o *TraceOptions) {
		o.IsPublicEndpoint = isPublic
	}
}

func WithSample(f func(c *gin.Context) trace.Sampler) TraceOption {
	return func(o *TraceOptions) {
		o.Sample = f
	}
}
