package frontends

import (
	"net/http"
	"time"
)

func ParseHttpTimespan(req *http.Request) (time.Time, time.Time, error) {
	intervalString := req.FormValue("interval")

	if intervalString != "" {
		// Parse stuff line "-1w1d"
		duration, err := parseDuration(intervalString)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		if duration > 0 {
			duration = -duration
		}
		return time.Now().Add(duration), time.Now(), nil
	}

	startString := req.FormValue("start")
	endString := req.FormValue("end")

	startDuration, err := parseDuration(startString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endDuration, err := parseDuration(endString)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return time.Now().Add(startDuration), time.Now().Add(endDuration), nil
}
