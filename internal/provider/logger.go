package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// tflogLogger provides a logging implementation using tflog.
type tflogLogger struct {
}

// Error logs an error message to tflog.
func (l *tflogLogger) Error(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Error(ctx, msg, additionalFields...)
}

// Warn logs a warning message to tflog.
func (l *tflogLogger) Warn(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Warn(ctx, msg, additionalFields...)
}

// Info logs an informational message to tflog.
func (l *tflogLogger) Info(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Info(ctx, msg, additionalFields...)
}

// Debug logs a debug message to tflog.
func (l *tflogLogger) Debug(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Debug(ctx, msg, additionalFields...)
}

// Trace logs a trace message to tflog.
func (l *tflogLogger) Trace(ctx context.Context, msg string, additionalFields ...map[string]any) {
	tflog.Trace(ctx, msg, additionalFields...)
}
