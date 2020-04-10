package ocgin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

type spanExporter struct {
	sync.Mutex
	cur []*trace.SpanData
}

var _ trace.Exporter = (*spanExporter)(nil)

func (se *spanExporter) ExportSpan(sd *trace.SpanData) {
	se.Lock()
	se.cur = append(se.cur, sd)
	se.Unlock()
}

func TestAlwaysSample(t *testing.T) {
	exporter := &spanExporter{cur: make([]*trace.SpanData, 0, 1)}
	trace.RegisterExporter(exporter)
	defer trace.UnregisterExporter(exporter)

	e := gin.Default()
	e.Use(HandlerFunc(WithSample(func(c *gin.Context) trace.Sampler {
		return trace.AlwaysSample()
	})))
	e.GET("/test", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	for i := 0; i <= 100; i++ {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			ctx := context.Background()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req = req.WithContext(ctx)
			e.ServeHTTP(w, req)
			cur := exporter.cur[i]
			if !cur.IsSampled() {
				t.Fatalf("SpanData is not sample")
			}
		})
	}
}

func TestSampleDecider(t *testing.T) {
	exporter := &spanExporter{cur: make([]*trace.SpanData, 0, 2)}
	trace.RegisterExporter(exporter)
	defer trace.UnregisterExporter(exporter)

	e := gin.Default()
	e.Use(HandlerFunc(WithSample(func(c *gin.Context) trace.Sampler {
		if c.FullPath() == "/health" {
			return trace.NeverSample()
		}
		return trace.AlwaysSample()
	})))
	e.GET("/test", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	e.GET("/health", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	tests := []struct {
		path       string
		wantSample bool
		message    string
	}{
		{path: "/test", wantSample: true, message: "sample success"},
		{path: "/health", wantSample: false, message: "not sample"},
	}

	for i, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			ctx := context.Background()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			req = req.WithContext(ctx)
			e.ServeHTTP(w, req)
			var got bool
			if len(exporter.cur) > i && exporter.cur[i].IsSampled() {
				got = true
			}
			if got != tt.wantSample {
				t.Fatalf("%s sample status wrong, want %v", tt.path, tt.wantSample)
			}
		})
	}

}
