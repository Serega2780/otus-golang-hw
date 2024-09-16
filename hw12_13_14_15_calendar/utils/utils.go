package utils

import "time"

func DayRange(date time.Time) (start time.Time, end time.Time) {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local),
		time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999, time.Local)
}

func WeekRange(date time.Time) (start time.Time, end time.Time) {
	tmp := date
	tmp = tmp.Add(7 * 24 * time.Hour)
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local),
		time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 23, 59, 59, 999, time.Local)
}

func MonthRange(date time.Time) (start time.Time, end time.Time) {
	tmp := date
	tmp = tmp.Add(30 * 24 * time.Hour)
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local),
		time.Date(tmp.Year(), tmp.Month(), tmp.Day(), 23, 59, 59, 999, time.Local)
}
