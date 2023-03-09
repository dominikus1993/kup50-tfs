package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatToAzureDevopsTime(t *testing.T) {
	day := time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local)

	subject := formatToAzureDevopsTime(day)

	assert.Equal(t, "05/01/2022", subject)
}

func TestFirstAndLastDayOfTheMonth(t *testing.T) {
	day := time.Date(2022, 5, 10, 0, 0, 0, 0, time.Local)

	first, last := FirstAndLastDayOfTheMonth(day)

	assert.Equal(t, "05/01/2022", first)
	assert.Equal(t, "05/31/2022", last)
}
