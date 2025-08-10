package helper

import (
	"log/slog"
	"strconv"
	"time"
)

func formatTime(i int64, unit string) string {
	if i == 1 {
		return "1 " + unit
	}
	slog.Info("formatTime", "value", i, "unit", unit)
	return strconv.FormatInt(i, 10) + " " + unit + "s"

}

func NowS() int64 {
	return time.Now().Unix()
}

func TimeFormat(timestamp int64) string {
	// 2023-10-01 12:00:00
	if timestamp <= 0 {
		slog.Warn("Invalid timestamp", "timestamp", timestamp)
		return "N/A"
	}
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
