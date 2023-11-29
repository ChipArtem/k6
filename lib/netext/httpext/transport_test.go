package httpext

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/ChipArtem/k6/lib"
	"github.com/ChipArtem/k6/metrics"
	"github.com/sirupsen/logrus"
)

func BenchmarkMeasureAndEmitMetrics(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	samples := make(chan metrics.SampleContainer, 10)
	defer close(samples)
	go func() {
		for range samples {
		}
	}()
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	registry := metrics.NewRegistry()
	state := &lib.State{
		Options: lib.Options{
			SystemTags: &metrics.DefaultSystemTagSet,
		},
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
		Samples:        samples,
		Logger:         logger,
	}
	t := transport{
		state:       state,
		ctx:         ctx,
		tagsAndMeta: &metrics.TagsAndMeta{Tags: registry.RootTagSet()},
	}

	unfRequest := &unfinishedRequest{
		tracer: &Tracer{},
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		request: &http.Request{
			URL: &url.URL{
				Host:   "example.com",
				Scheme: "https",
			},
		},
	}

	b.Run("no responseCallback", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t.measureAndEmitMetrics(unfRequest)
		}
	})

	t.responseCallback = func(n int) bool { return true }

	b.Run("responseCallback", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			t.measureAndEmitMetrics(unfRequest)
		}
	})
}
