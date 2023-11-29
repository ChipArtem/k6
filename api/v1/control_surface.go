package v1

import (
	"context"

	"github.com/ChipArtem/k6/execution"
	"github.com/ChipArtem/k6/lib"
	"github.com/ChipArtem/k6/metrics"
	"github.com/ChipArtem/k6/metrics/engine"
)

// ControlSurface includes the methods the REST API can use to control and
// communicate with the rest of k6.
type ControlSurface struct {
	RunCtx        context.Context
	Samples       chan metrics.SampleContainer
	MetricsEngine *engine.MetricsEngine
	Scheduler     *execution.Scheduler
	RunState      *lib.TestRunState
}
