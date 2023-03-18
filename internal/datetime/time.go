package datetime

import "time"

// 6/14/2018
func FirstAndLastDayOfTheMonth(now time.Time) (first time.Time, last time.Time) {
	year, month, _ := now.Date()
	currentLocation := now.Location()
	firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	first = firstDayOfMonth
	last = lastOfMonth
	return
}

func FormatToAzureDevopsTime(t time.Time) *string {
	result := t.Format("01/02/2006")
	return &result
}
