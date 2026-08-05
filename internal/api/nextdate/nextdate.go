package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	errRepeatIsEmpty       = errors.New("repeat is empty")
	errWrongSymbol         = errors.New("wrong symbol")
	errIntervalIsNotSet    = errors.New("interval is not set")
	errWrongInterval       = errors.New("wrong interval")
	errUnsupportedFormat   = errors.New("unsupported format")
	errExceededMaxInterval = errors.New("the maximum allowed interval has been exceeded")
)

const dateFormat string = "20060102"

type interval struct {
	days   int
	months int
	years  int
}

func afterNow(date time.Time, now time.Time) bool {
	return date.After(now)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", errRepeatIsEmpty
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	repeatSlice := strings.Split(repeat, " ")
	if len(repeatSlice) == 1 && repeatSlice[0] != "y" {
		return "", errIntervalIsNotSet
	}

	var i interval

	switch repeatSlice[0] {
	case "y":
		if len(repeatSlice) > 1 {
			return "", errWrongInterval
		}
		i.years = 1
	case "d":
		if len(repeatSlice) > 2 {
			return "", errWrongInterval
		}

		i.days, err = strconv.Atoi(repeatSlice[1])
		if err != nil {
			return "", err
		}
		if i.days > 400 {
			return "", errExceededMaxInterval
		}
	case "w":
		if len(repeatSlice) > 2 {
			return "", errWrongInterval
		}

		//weekdaySlice := strings.Split(repeatSlice[1], ",")

		return "", errUnsupportedFormat
	case "m":
		return "", errUnsupportedFormat
	default:
		return "", errWrongSymbol
	}

	for {
		date = date.AddDate(i.years, i.months, i.days)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}
