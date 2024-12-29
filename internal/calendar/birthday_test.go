package calendar

import (
	"testing"
	"time"
)

func TestBirthdayDateBefore(t *testing.T) {
	testcases := []struct {
		title      string
		day        uint8
		month      time.Month
		otherDay   uint8
		otherMonth time.Month
		isBefore   bool
	}{
		{
			title:      "before",
			day:        1,
			month:      time.February,
			otherDay:   2,
			otherMonth: time.April,
			isBefore:   true,
		},
		{
			title:      "before, same month",
			day:        1,
			month:      time.May,
			otherDay:   2,
			otherMonth: time.May,
			isBefore:   true,
		},
		{
			title:      "before, same day",
			day:        1,
			month:      time.February,
			otherDay:   1,
			otherMonth: time.April,
			isBefore:   true,
		},
		{
			title:      "after, same date",
			day:        1,
			month:      time.January,
			otherDay:   1,
			otherMonth: time.January,
			isBefore:   false,
		},
		{
			title:      "after",
			day:        19,
			month:      time.May,
			otherDay:   23,
			otherMonth: time.January,
			isBefore:   false,
		},
	}

	for _, test := range testcases {
		t.Run(test.title, func(t *testing.T) {
			bDate := BirthdayDate{Month: test.month, Day: test.day}
			otherDate := BirthdayDate{Month: test.otherMonth, Day: test.otherDay}

			isBefore := bDate.Before(otherDate)
			if test.isBefore != isBefore {
				t.Errorf("Before(): expected %v got %v", test.isBefore, isBefore)
			}
		})
	}
}
