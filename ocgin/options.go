package ocgin

import (
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

type TraceOption func(o *TraceOptions)

type TraceOptions struct {
	// defaultAttributes will be set to each span as default.
	defaultAttributes []trace.Attribute

	// isPublicEndpoint should be set to true for publicly accessible HTTP(S)
	// servers. If true, any trace metadata set on the incoming request will
	// be added as a linked trace instead of being added as a parent of the
	// current trace.
	isPublicEndpoint bool

	// sample is applied to the span started by this Handler around each
	// request. default is 0.0001
	sample func(c *gin.Context) trace.Sampler
}

var _defaultOptions = TraceOptions{
	defaultAttributes: []trace.Attribute{},
	isPublicEndpoint:  false,
}

// WithDefaultAttributes will be set to each span as default.
func WithDefaultAttributes(attrs ...trace.Attribute) TraceOption {
	return func(o *TraceOptions) {
		o.defaultAttributes = attrs
	}
}

// WithPublicEndpoint receive true when server is public
func WithPublicEndpoint(isPublic bool) TraceOption {
	return func(o *TraceOptions) {
		o.isPublicEndpoint = isPublic
	}
}

// WithSample receive a function with gin.Context to decide sample with each request
func WithSample(f func(c *gin.Context) trace.Sampler) TraceOption {
	return func(o *TraceOptions) {
		o.sample = f
	}
}
