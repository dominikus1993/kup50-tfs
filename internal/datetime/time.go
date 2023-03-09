package datetime

import "time"

// 6/14/2018
func FirstAndLastDayOfTheMonth(now time.Time) (first string, last string) {
	year, month, _ := now.Date()
	currentLocation := now.Location()
	firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	first = formatToAzureDevopsTime(firstDayOfMonth)
	last = formatToAzureDevopsTime(lastOfMonth)
	return
}

func formatToAzureDevopsTime(t time.Time) string {
	return t.Format("01/02/2006")
}
