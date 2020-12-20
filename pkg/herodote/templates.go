package herodote

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ViBiOh/herodote/pkg/model"
)

const (
	daysInWeek   = float64(7)
	weeksInMonth = float64(4)
	monthsInYear = float64(12)
)

var (
	repositoriesColors = []string{
		"#006E6D",
		"#2A4B7C",
		"#3F69AA",
		"#77212E",
		"#577284",
		"#6C4F3D",
		"#797B3A",
		"#935529",
		"#BD3D3A",
		"#9B1B30",
		"#E08119",
		"#6B5B95",
		"#F96714",
		"#485167",
		"#2E4A62",
		"#264E36",
	}
	colorsCount = 0
	colors      = sync.Map{}

	// FuncMap for template rendering
	FuncMap = template.FuncMap{
		"colors": func(commit model.Commit) string {
			if color, ok := colors.Load(commit.Repository); ok {
				return color.(string)
			}

			colorsCount++
			nextColor := repositoriesColors[colorsCount%len(repositoriesColors)]
			colors.Store(commit.Repository, nextColor)

			return nextColor
		},
		"contains":           contains,
		"dateDistanceInDays": diffInDays,
		"toggleParam": func(path string, params url.Values, name, value string) string {
			safeValues := url.Values{}
			done := false

			for key := range params {
				currentValue := strings.TrimSpace(params.Get(key))

				if len(currentValue) == 0 {
					continue
				}

				if key == name {
					done = true
				} else {
					safeValues.Set(key, currentValue)
				}
			}

			if !done {
				safeValues.Set(name, value)
			}

			if len(safeValues) == 0 {
				return path
			}

			return fmt.Sprintf("%s?%s", path, safeValues.Encode())
		},
	}
)

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if strings.EqualFold(item, value) {
			return true
		}
	}

	return false
}

func diffInDays(date, now time.Time) string {
	beginNow := now.Truncate(dayDuration)
	beginDate := date.Truncate(dayDuration)

	if beginNow.Unix() == beginDate.Unix() {
		return "Today"
	}

	count := math.Abs(beginNow.Sub(beginDate).Truncate(dayDuration).Hours()) / 24

	if count < daysInWeek {
		if count < 2 {
			return "Yesterday"
		}

		return fmt.Sprintf("%.f days ago", count)
	}

	count = count / daysInWeek
	if count < weeksInMonth {
		if count < 2 {
			return "1 week ago"
		}

		return fmt.Sprintf("%.f weeks ago", count)
	}

	count = count / weeksInMonth
	if count < monthsInYear {
		if count < 2 {
			return "1 month ago"
		}

		return fmt.Sprintf("%.f months ago", count)
	}

	count = count / monthsInYear
	if count < 2 {
		return "1 year ago"
	}

	return fmt.Sprintf("%.f years ago", count)
}
