package slogs

import "log/slog"

// CLog returns a child logger.
func CLog(subsys string) *slog.Logger {
	return slog.With(Subsys, subsys)
}
