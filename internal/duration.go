package internal

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var durationRegex = regexp.MustCompile("([0-9]+[ywdhm])")

func ParseDuration(s string) (time.Duration, error) {
	matches := durationRegex.FindAllStringSubmatch(s, -1)
	d := 0

	found := ""
	for _, match := range matches {
		match := match[1]

		// check unit
		unit := string(match[len(match)-1])

		if strings.Contains(found, unit) {
			return time.Duration(0), errors.New("invalid format")
		}
		// check if we haven't seen it already
		found = found + unit

		// Check if quantity is valid
		qty, err := strconv.Atoi(match[:len(match)-1])
		if err != nil {
			return time.Duration(0), err
		}

		switch unit {
		case "y": // year
			d += qty * 365 * 24 * 3600
		case "w": // week
			d += qty * 7 * 24 * 3600
		case "d": // day
			d += qty * 24 * 3600
		case "h": // hour
			d += qty * 3600
		case "m": // minute
			d += qty * 60
		}
	}

	return time.Duration(d * int(time.Second)), nil
}
