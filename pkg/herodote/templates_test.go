package herodote

import (
	"testing"
	"time"
)

func TestDiffInDays(t *testing.T) {
	type args struct {
		date time.Time
		now  time.Time
	}

	var cases = []struct {
		intention string
		args      args
		want      string
	}{
		{
			"Same day",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 20, 6, 45, 0, 0, time.UTC),
			},
			"Today",
		},
		{
			"Yesterday",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 19, 6, 45, 0, 0, time.UTC),
			},
			"Yesterday",
		},
		{
			"Some days",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 17, 6, 45, 0, 0, time.UTC),
			},
			"3 days ago",
		},
		{
			"One week",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 13, 6, 45, 0, 0, time.UTC),
			},
			"1 week ago",
		},
		{
			"Some weeks",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 1, 6, 45, 0, 0, time.UTC),
			},
			"3 weeks ago",
		},
		{
			"One month",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 11, 20, 6, 45, 0, 0, time.UTC),
			},
			"1 month ago",
		},
		{
			"Some months",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 9, 20, 6, 45, 0, 0, time.UTC),
			},
			"3 months ago",
		},
		{
			"One year",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2019, 9, 20, 6, 45, 0, 0, time.UTC),
			},
			"1 year ago",
		},
		{
			"Some years",
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2017, 9, 20, 6, 45, 0, 0, time.UTC),
			},
			"4 years ago",
		},
	}

	for _, tc := range cases {
		t.Run(tc.intention, func(t *testing.T) {
			if got := diffInDays(tc.args.date, tc.args.now); got != tc.want {
				t.Errorf("diffInDays() = `%s`, want `%s`", got, tc.want)
			}
		})
	}
}
