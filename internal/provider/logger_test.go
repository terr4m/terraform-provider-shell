package provider

import (
	"testing"
)

func Test_tflogLogger(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	logger := tflogLogger{}
	logger.Error(ctx, "error")
	logger.Warn(ctx, "warn")
	logger.Info(ctx, "info")
	logger.Debug(ctx, "debug")
	logger.Trace(ctx, "trace")
}
