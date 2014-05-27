package frontends

import (
	"net/http"
	"strconv"
	"time"
)

// Parse an interval and give start and end estimates
func parseInterval(interval string) (time.Time, time.Time, error) {
	// If it's all digits - YYYYMMDD, YYYYMM or YYYY
	_, err := strconv.ParseInt(interval, 10, 32)
	if err == nil && (len(interval) == 8 || len(interval) == 6 || len(interval) == 4) {
		startMonth := time.January
		endMonth := time.December
		startDay := 1
		endDay := 31

		// There's always a year
		year, _ := strconv.ParseInt(interval[0:4], 10, 32)
		// endDay = 31 // Always ends on december 31st

		if len(interval) >= 6 {
			month, _ := strconv.ParseInt(interval[4:6], 10, 32)
			startMonth = time.Month(month)
			endMonth = time.Month(month)
			// TODO: Figure out last day of last month
		}
		if len(interval) == 8 {
			day, _ := strconv.ParseInt(interval[6:8], 10, 32)
			startDay = int(day)
			endDay = int(day)
		}

		return time.Date(int(year), startMonth, startDay, 0, 0, 0, 0, time.UTC), time.Date(int(year), endMonth, endDay, 23, 59, 59, 0, time.UTC), nil
	}

	// Relative - ex. -1w
	duration, err := parseDuration(interval)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Flip duration so it always goes back in time
	if duration > 0 {
		duration = -duration
	}

	return time.Now().In(time.UTC).Add(duration), time.Now().In(time.UTC), nil
}

func ParseHttpTimespan(req *http.Request) (time.Time, time.Time, error) {
	intervalString := req.FormValue("interval")

	if intervalString != "" {
		return parseInterval(intervalString)
	}

	startString := req.FormValue("start")
	endString := req.FormValue("end")

	startTime, _, err := parseInterval(startString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	_, endTime, err := parseInterval(endString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if endTime.Before(startTime) {
		endTime, startTime = startTime, endTime
	}

	return startTime, endTime, nil
}
