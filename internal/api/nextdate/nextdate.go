package nextdate

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	errRepeatIsEmpty       = errors.New("repeat is empty")
	errWrongSymbol         = errors.New("wrong symbol")
	errIntervalIsNotSet    = errors.New("interval is not set")
	errInvalidInterval     = errors.New("invalid interval")
	errUnsupportedFormat   = errors.New("unsupported format")
	errExceededMaxInterval = errors.New("the maximum allowed interval has been exceeded")
	errInvalidValue        = errors.New("invalid value")
	errInvalidDayOfMonth   = errors.New("invalid day of the month")
	errInvalidMonth        = errors.New("invalid month")
)

const dateFormat string = "20060102"

type interval struct {
	days   int
	months int
	years  int
}

func AfterNow(date time.Time, now time.Time) bool {
	return date.Format("20060102") > now.Format("20060102")
}

type weekdaysRules struct {
	rule []int
}

func (w *weekdaysRules) parseWeekdayRules(repeatStr string) error {
	weekdaySlice := strings.Split(repeatStr, ",")
	for _, weekday := range weekdaySlice {
		weekdayInt, err := strconv.Atoi(weekday)
		if err != nil {
			return err
		}
		if weekdayInt < 1 || weekdayInt > 7 {
			return errInvalidValue
		}
		w.rule = append(w.rule, weekdayInt)
	}

	sort.Ints(w.rule) // для сортировки слайса дней недели, если они не по порядку

	return nil
}

func (w *weekdaysRules) countDays(weekdayNow int, date time.Time, now time.Time) int {
	days := 0
	for idx, weekday := range w.rule {
		if weekday > weekdayNow {
			days = weekday - weekdayNow
			break
		} else if idx == len(w.rule)-1 {
			dur := now.Sub(date)
			dayDur := int(dur / (24 * time.Hour))
			days = dayDur + (7 - weekdayNow + w.rule[0])
		}
	}
	return days
}

type daysOfTheMonthsRules struct {
	daysRule       [32]bool
	monthsRule     [13]bool
	lastDayRule    bool
	preLastDayRule bool
}

func (d *daysOfTheMonthsRules) parseDaysOfTheMonthsRules(repeatSlice []string) error {
	if len(repeatSlice) == 2 {
		for month := 1; month < len(d.monthsRule); month++ {
			d.monthsRule[month] = true
		}
	}

	for i := 1; i < len(repeatSlice); i++ {
		infoStrSlice := strings.Split(repeatSlice[i], ",")
		for _, infoStr := range infoStrSlice {
			info, err := strconv.Atoi(infoStr)
			if err != nil {
				return err
			}

			switch i {
			case 1:
				if info < -2 || info > 31 || info == 0 {
					return errInvalidDayOfMonth
				}
				switch info {
				case -2:
					d.preLastDayRule = true
				case -1:
					d.lastDayRule = true
				default:
					d.daysRule[info] = true
				}
			case 2:
				if info < 1 || info > 12 {
					return errInvalidMonth
				}
				d.monthsRule[info] = true
			}
		}
	}
	return nil
}

func (d *daysOfTheMonthsRules) countDays(date time.Time, now time.Time) int {
	dateTmp := date
	daysCount := 0
	for {
		daysCount++
		dateTmp = dateTmp.AddDate(0, 0, 1)
		year, month, day := dateTmp.Date()

		isDaysRule := d.daysRule[day]
		isMonthRule := d.monthsRule[month]

		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, dateTmp.Location())
		preLastDay := lastDay.Day() - 1

		if d.preLastDayRule && day == preLastDay {
			isDaysRule = true
		}
		if d.lastDayRule && day == lastDay.Day() {
			isDaysRule = true
		}

		if AfterNow(dateTmp, now) && isDaysRule && isMonthRule {
			break
		}
	}
	return daysCount
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
			return "", errInvalidInterval
		}

		i.years = 1

	case "d":
		if len(repeatSlice) > 2 {
			return "", errInvalidInterval
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
			return "", errInvalidInterval
		}

		var wRules weekdaysRules

		err = wRules.parseWeekdayRules(repeatSlice[1])
		if err != nil {
			return "", err
		}

		weekdayNow := int(now.Weekday())

		i.days = wRules.countDays(weekdayNow, date, now)

	case "m":
		if len(repeatSlice) > 3 {
			return "", errInvalidInterval
		}

		var dmRules daysOfTheMonthsRules

		err := dmRules.parseDaysOfTheMonthsRules(repeatSlice)
		if err != nil {
			return "", err
		}

		i.days = dmRules.countDays(date, now)

	default:
		return "", errWrongSymbol
	}

	for {
		date = date.AddDate(i.years, i.months, i.days)
		if AfterNow(date, now) {
			break
		}
	}

	return date.Format(dateFormat), nil
}
