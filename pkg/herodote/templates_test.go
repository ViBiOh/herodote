package herodote

import (
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	type args struct {
		arr   []string
		value string
	}

	cases := map[string]struct {
		args args
		want bool
	}{
		"nil": {
			args{
				arr:   nil,
				value: "hello",
			},
			false,
		},
		"simple": {
			args{
				arr:   []string{"hello", "world"},
				value: "hello",
			},
			true,
		},
		"absent": {
			args{
				arr:   []string{"world"},
				value: "hello",
			},
			false,
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			if got := contains(tc.args.arr, tc.args.value); got != tc.want {
				t.Errorf("contains() = %t, want %t", got, tc.want)
			}
		})
	}
}

func TestDiffInDays(t *testing.T) {
	type args struct {
		date time.Time
		now  time.Time
	}

	cases := map[string]struct {
		args args
		want string
	}{
		"Same day": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 20, 6, 45, 0, 0, time.UTC),
			},
			"Today",
		},
		"Yesterday": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 19, 6, 45, 0, 0, time.UTC),
			},
			"Yesterday",
		},
		"Some days": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 17, 6, 45, 0, 0, time.UTC),
			},
			"3 days ago",
		},
		"One week": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 13, 6, 45, 0, 0, time.UTC),
			},
			"1 week ago",
		},
		"Some weeks": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 12, 1, 6, 45, 0, 0, time.UTC),
			},
			"3 weeks ago",
		},
		"One month": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 11, 20, 6, 45, 0, 0, time.UTC),
			},
			"1 month ago",
		},
		"Some months": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2020, 9, 20, 6, 45, 0, 0, time.UTC),
			},
			"3 months ago",
		},
		"One year": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2019, 9, 20, 6, 45, 0, 0, time.UTC),
			},
			"1 year ago",
		},
		"Some years": {
			args{
				now:  time.Date(2020, 12, 20, 18, 45, 0, 0, time.UTC),
				date: time.Date(2017, 9, 20, 6, 45, 0, 0, time.UTC),
			},
			"4 years ago",
		},
	}

	for intention, tc := range cases {
		t.Run(intention, func(t *testing.T) {
			if got := diffInDays(tc.args.date, tc.args.now); got != tc.want {
				t.Errorf("diffInDays() = `%s`, want `%s`", got, tc.want)
			}
		})
	}
}
