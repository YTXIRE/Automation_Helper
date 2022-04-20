package times

import (
	"time"
)

func ConvertTime(countTime int64, timeType string) time.Duration {
	switch timeType {
	case "s":
		return time.Duration(countTime) * time.Second
	case "m":
		return time.Duration(countTime) * time.Minute
	case "h":
		return time.Duration(countTime) * time.Hour
	case "d":
		return time.Duration(countTime) * 24 * time.Hour
	case "w":
		return time.Duration(countTime) * 7 * 24 * time.Hour
	case "M":
		return time.Duration(countTime) * 4 * 7 * 24 * time.Hour
	case "Y":
		return time.Duration(countTime) * 12 * 4 * 7 * 24 * time.Hour
	}
	return time.Duration(countTime) * time.Hour
}
