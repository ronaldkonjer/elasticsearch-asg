package health

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/heptiolabs/healthcheck"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/mintel/elasticsearch-asg/mocks/mockhttp"
)

func setup(t *testing.T, checkFactory func(context.Context, string) healthcheck.Check) (healthcheck.Check, *httptest.Server, *mockhttp.Mux, func()) {
	logger := zaptest.NewLogger(t)
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()
	t1 := zap.ReplaceGlobals(logger)
	t2 := zap.RedirectStdLog(logger)

	ctx, cancel := context.WithCancel(context.Background())
	server, mux := mockhttp.NewServer()
	check := checkFactory(ctx, server.URL)

	originalTimeout := DefaultHTTPTimeout
	DefaultHTTPTimeout = 500 * time.Millisecond

	return check, server, mux, func() {
		DefaultHTTPTimeout = originalTimeout
		cancel()
		server.Close()
		t2()
		t1()
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}
}
