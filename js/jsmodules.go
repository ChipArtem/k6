package js

import (
	"github.com/ChipArtem/k6/ext"
	"github.com/ChipArtem/k6/js/modules/k6"
	"github.com/ChipArtem/k6/js/modules/k6/crypto"
	"github.com/ChipArtem/k6/js/modules/k6/crypto/x509"
	"github.com/ChipArtem/k6/js/modules/k6/data"
	"github.com/ChipArtem/k6/js/modules/k6/encoding"
	"github.com/ChipArtem/k6/js/modules/k6/execution"
	"github.com/ChipArtem/k6/js/modules/k6/experimental/fs"
	"github.com/ChipArtem/k6/js/modules/k6/experimental/tracing"
	"github.com/ChipArtem/k6/js/modules/k6/grpc"
	"github.com/ChipArtem/k6/js/modules/k6/html"
	"github.com/ChipArtem/k6/js/modules/k6/http"
	"github.com/ChipArtem/k6/js/modules/k6/metrics"
	"github.com/ChipArtem/k6/js/modules/k6/ws"

	"github.com/grafana/xk6-browser/browser"
	expGrpc "github.com/grafana/xk6-grpc/grpc"
	"github.com/grafana/xk6-redis/redis"
	"github.com/grafana/xk6-timers/timers"
	"github.com/grafana/xk6-webcrypto/webcrypto"
	expws "github.com/grafana/xk6-websockets/websockets"
)

func getInternalJSModules() map[string]interface{} {
	return map[string]interface{}{
		"k6":                         k6.New(),
		"k6/crypto":                  crypto.New(),
		"k6/crypto/x509":             x509.New(),
		"k6/data":                    data.New(),
		"k6/encoding":                encoding.New(),
		"k6/execution":               execution.New(),
		"k6/experimental/redis":      redis.New(),
		"k6/experimental/webcrypto":  webcrypto.New(),
		"k6/experimental/websockets": &expws.RootModule{},
		"k6/experimental/grpc":       expGrpc.New(),
		"k6/experimental/timers":     timers.New(),
		"k6/experimental/tracing":    tracing.New(),
		"k6/experimental/browser":    browser.New(),
		"k6/experimental/fs":         fs.New(),
		"k6/net/grpc":                grpc.New(),
		"k6/html":                    html.New(),
		"k6/http":                    http.New(),
		"k6/metrics":                 metrics.New(),
		"k6/ws":                      ws.New(),
	}
}

func getJSModules() map[string]interface{} {
	result := getInternalJSModules()
	external := ext.Get(ext.JSExtension)

	// external is always prefixed with `k6/x`
	for _, e := range external {
		result[e.Name] = e.Module
	}

	return result
}
