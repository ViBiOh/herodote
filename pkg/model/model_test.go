package model

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSanitize(t *testing.T) {
	var cases = []struct {
		intention string
		instance  Commit
		want      Commit
	}{
		{
			"simple",
			Commit{
				Hash:       "  Hash   ",
				Type:       "  Type   ",
				Component:  "  Component   ",
				Content:    "  Content   ",
				Remote:     "  Remote   ",
				Repository: "  Repository   ",
			},
			Commit{
				Hash:       "hash",
				Type:       "type",
				Component:  "component",
				Content:    "  Content   ",
				Remote:     "remote",
				Repository: "repository",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.intention, func(t *testing.T) {
			if got := tc.instance.Sanitize(); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Sanitize() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	var cases = []struct {
		intention string
		instance  Commit
		wantErr   error
	}{
		{
			"empty",
			Commit{},
			errors.New("commit's hash is required"),
		},
		{
			"type",
			Commit{
				Hash: "1ab2c3f4d",
			},
			errors.New("commit's type is required"),
		},
		{
			"content",
			Commit{
				Hash: "1ab2c3f4d",
				Type: "feat",
			},
			errors.New("commit's content is required"),
		},
		{
			"date",
			Commit{
				Hash:    "1ab2c3f4d",
				Type:    "feat",
				Content: "Add README.md",
			},
			errors.New("commit's date is required"),
		},
		{
			"remote",
			Commit{
				Hash:    "1ab2c3f4d",
				Type:    "feat",
				Content: "Add CONTRIBUTING.md",
				Date:    time.Now(),
			},
			errors.New("repository's remote is required"),
		},
		{
			"repository",
			Commit{
				Hash:    "1ab2c3f4d",
				Type:    "feat",
				Content: "Add LICENSE",
				Date:    time.Now(),
				Remote:  "github.com",
			},
			errors.New("repository's name is required"),
		},
		{
			"valid",
			Commit{
				Hash:       "1ab2c3f4d",
				Type:       "feat",
				Content:    "Add main.go",
				Date:       time.Now(),
				Remote:     "github.com",
				Repository: "vibioh/herodote",
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.intention, func(t *testing.T) {
			gotErr := tc.instance.Check()

			failed := false

			if tc.wantErr == nil && gotErr != nil {
				failed = true
			} else if tc.wantErr != nil && gotErr == nil {
				failed = true
			} else if tc.wantErr != nil && !strings.Contains(gotErr.Error(), tc.wantErr.Error()) {
				failed = true
			}

			if failed {
				t.Errorf("Check() = `%s`, want `%s`", gotErr, tc.wantErr)
			}
		})
	}
}
