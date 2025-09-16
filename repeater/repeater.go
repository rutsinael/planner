package repeater

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func GetNextDate(now time.Time, dstart string, repeat string) (string, error) {

	date, err := time.Parse(DateFormat, dstart)

	if err != nil {
		return "", err
	}

	if len(repeat) == 0 {
		return "", errors.New("no spitted rule")
	}
	spitedRule := strings.Split(repeat, " ")

	var validationError error

	switch spitedRule[0] {
	case "d":
		date, validationError = countNextFromDay(now, date, repeat)
	case "y":
		date, validationError = countNextFromYear(now, date, repeat)
	case "w":
		date, validationError = countNextFromWeek(now, date, repeat)
	case "m":
		date, validationError = countNextFromMonth(now, date, repeat)
	default:
		return "", errors.New("unknown rule")
	}

	convertedDate := date.Format(DateFormat)
	return convertedDate, validationError
}

func AfterNow(date, now time.Time) bool {

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.After(now)
}

func countNextFromYear(now, date time.Time, rule string) (time.Time, error) {

	if rule != "y" {
		return time.Time{}, errors.New("invalid year rule")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if AfterNow(date, now) {
			break
		}
	}

	return date, nil
}

func countNextFromDay(now, date time.Time, rule string) (time.Time, error) {

	rules := strings.Split(rule, " ")

	if len(rules) != 2 {
		return time.Time{}, errors.New("invalid day rule")
	}

	interval, err := strconv.Atoi(rules[1])

	if err != nil {
		return time.Time{}, errors.New("invalid day interval")
	}

	if interval > 400 {
		return time.Time{}, errors.New("invalid day interval: day couldn't be bigger than 400")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if AfterNow(date, now) {
			break
		}
	}

	return date, nil
}

func countNextFromWeek(now, date time.Time, rule string) (time.Time, error) {

	rules := strings.Split(rule, " ")

	if len(rules) != 2 {
		return time.Time{}, errors.New("invalid week rule")
	}

	intervals := strings.Split(rules[1], ",")

	weekDays := make([]int, 0)
	for _, day := range intervals {
		j, err := strconv.Atoi(day)

		if err != nil || j > 7 || j < 1 {
			return time.Time{}, errors.New("invalid week interval")
		}

		if j == 7 {
			weekDays = append(weekDays, 0)
		} else {
			weekDays = append(weekDays, j)
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if AfterNow(date, now) && slices.Contains(weekDays, int(date.Weekday())) {
			break
		}
	}

	return date, nil
}

func countNextFromMonth(now, date time.Time, rule string) (time.Time, error) {

	rules := strings.Split(rule, " ")
	if len(rules) < 2 {
		return time.Time{}, errors.New("invalid month rule")
	}

	day, month := initDaysAndMonths()

	daysRule := strings.Split(rules[1], ",")
	for _, dr := range daysRule {
		d, err := strconv.Atoi(dr)

		if err != nil {
			return time.Time{}, errors.New("invalid month interval")
		}

		if _, ok := day[d]; !ok {
			return time.Time{}, errors.New("invalid day")
		}
		day[d] = true
	}

	if len(rules) == 3 {
		monthRule := strings.Split(rules[2], ",")
		for _, mr := range monthRule {
			m, err := strconv.Atoi(mr)

			if err != nil {
				return time.Time{}, errors.New("invalid month interval")
			}

			_, ok := month[m]
			if !ok {
				return time.Time{}, errors.New("invalid month")
			}
			month[m] = true
		}
	}

	for {
		date = date.AddDate(0, 0, 1)

		md := daysInMonth(date)
		okDay := day[date.Day()] ||
			(day[-1] && date.Day() == md) ||
			(day[-2] && date.Day() == md-1)

		okMonth := len(rules) == 2 || month[int(date.Month())]

		if AfterNow(date, now) && okDay && okMonth {
			break
		}

	}

	return date, nil
}

func initDaysAndMonths() (map[int]bool, map[int]bool) {
	day := make(map[int]bool)
	for i := 1; i < 32; i++ {
		day[i] = false
	}
	day[-1] = false
	day[-2] = false

	month := make(map[int]bool)
	for i := 1; i < 13; i++ {
		month[i] = false
	}

	return day, month
}

func daysInMonth(t time.Time) int {
	firstDayOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	lastDayOfCurrentMonth := firstDayOfNextMonth.AddDate(0, 0, -1)
	return lastDayOfCurrentMonth.Day()
}
