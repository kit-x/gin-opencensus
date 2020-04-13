package ocgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

var _defaultFormat propagation.HTTPFormat = &b3.HTTPFormat{}

func HandlerFunc(opts ...TraceOption) gin.HandlerFunc {
	opt := _defaultOptions
	for _, f := range opts {
		f(&opt)
	}
	return func(c *gin.Context) {
		if c.FullPath() == "" {
			c.Next()
			return
		}
		ctx := c.Request.Context()
		name := formatSpanName(c)

		startOptions := append(
			make([]trace.StartOption, 0, 2),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithSampler(opt.sample(c)),
		)

		// Code reference https://github.com/census-instrumentation/opencensus-go/blob/master/plugin/ochttp/server.go
		var span *trace.Span
		sc, ok := _defaultFormat.SpanContextFromRequest(c.Request)
		if ok && !opt.isPublicEndpoint {
			ctx, span = trace.StartSpanWithRemoteParent(ctx, name, sc, startOptions...)
		} else {
			ctx, span = trace.StartSpan(c.Request.Context(), name, startOptions...)

			if ok {
				span.AddLink(trace.Link{
					TraceID: sc.TraceID,
					SpanID:  sc.SpanID,
					Type:    trace.LinkTypeParent,
				})
			}
		}
		defer span.End()

		attrs := append(requestAttrs(c), opt.defaultAttributes...)
		span.AddAttributes(attrs...)
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			span.AddMessageReceiveEvent(0,
				c.Request.ContentLength, -1)
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		c.Header(b3.TraceIDHeader, span.SpanContext().TraceID.String())

		statusCode := c.Writer.Status()
		span.AddAttributes(trace.Int64Attribute(ochttp.StatusCodeAttribute, int64(statusCode)))
		span.SetStatus(ochttp.TraceStatus(statusCode, http.StatusText(statusCode)))
	}
}

func formatSpanName(c *gin.Context) string {
	return c.Request.Method + " " + c.FullPath()
}

func requestAttrs(c *gin.Context) []trace.Attribute {
	attrs := make([]trace.Attribute, 0, 5)
	attrs = append(attrs,
		trace.StringAttribute(ochttp.PathAttribute, c.FullPath()),
		trace.StringAttribute(ochttp.URLAttribute, c.Request.URL.String()),
		trace.StringAttribute(ochttp.HostAttribute, c.Request.Host),
		trace.StringAttribute(ochttp.MethodAttribute, c.Request.Method),
	)

	userAgent := c.Request.UserAgent()
	if userAgent != "" {
		attrs = append(attrs, trace.StringAttribute(ochttp.UserAgentAttribute, userAgent))
	}

	return attrs
}
