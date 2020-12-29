package internal

import "time"

// TimeNow returns the current time in UTC zone WITHOUT nanoseconds. This is useful when marshalling times to JSON
func TimeNow() time.Time {
	ct := time.Now().Unix()
	return time.Unix(ct, 0).UTC()
}
