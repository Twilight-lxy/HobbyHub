package utils

import (
	"time"
)

func GetCurrentTime() time.Time {
	return time.Now()
}
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
func GetCurrentTimestampMillis() int64 {
	return time.Now().UnixMilli()
}
func GetCurrentTimestampNanos() int64 {
	return time.Now().UnixNano()
}
func GetCurrentTimeFormatted() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func GetCurrentTimeFormattedWithTimezone() string {
	return time.Now().Format("2006-01-02 15:04:05 -0700")
}
func ParseTimeFromString(timeStr string) time.Time {
	timeStrm, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Time{}
	}
	return timeStrm
}
func ParseTimeFromStringWithTimezone(timeStr string) time.Time {
	timeStrm, err := time.Parse("2006-01-02 15:04:05 -0700", timeStr)
	if err != nil {
		return time.Time{}
	}
	return timeStrm
}
func FormatTimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
func FormatTimeToStringWithTimezone(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 -0700")
}
func TimeToUnixTimestamp(t time.Time) int64 {
	return t.Unix()
}
func TimeToUnixTimestampMillis(t time.Time) int64 {
	return t.UnixMilli()
}
func TimeToUnixTimestampNanos(t time.Time) int64 {
	return t.UnixNano()
}
