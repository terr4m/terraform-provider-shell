package script

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TFLogLogger provides a logging implementation using tflog.
type TFLogLogger struct{}

// Error logs an error message to tflog.
func (l *TFLogLogger) Error(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Error(ctx, msg, additionalFields...)
}

// Warn logs a warning message to tflog.
func (l *TFLogLogger) Warn(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Warn(ctx, msg, additionalFields...)
}

// Info logs an informational message to tflog.
func (l *TFLogLogger) Info(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Info(ctx, msg, additionalFields...)
}

// Debug logs a debug message to tflog.
func (l *TFLogLogger) Debug(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Debug(ctx, msg, additionalFields...)
}

// Trace logs a trace message to tflog.
func (l *TFLogLogger) Trace(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Trace(ctx, msg, additionalFields...)
}
