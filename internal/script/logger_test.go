package script_test

import (
	"testing"

	"github.com/terr4m/terraform-provider-shell/internal/script"
)

func TestTFLogLogger(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	logger := script.TFLogLogger{}
	logger.Error(ctx, "error")
	logger.Warn(ctx, "warn")
	logger.Info(ctx, "info")
	logger.Debug(ctx, "debug")
	logger.Trace(ctx, "trace")
}
