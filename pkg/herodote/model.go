package herodote

import (
	"fmt"
	"strings"
)

// Commit describes a commit on a repository
type Commit struct {
	Hash       string `json:"hash"`
	Type       string `json:"type"`
	Component  string `json:"component"`
	Revert     bool   `json:"revert"`
	Breaking   bool   `json:"breaking"`
	Content    string `json:"content"`
	Date       uint64 `json:"date"`
	Remote     string `json:"remote"`
	Repository string `json:"repository"`
}

// Sanitize cleans value of a commit
func (c Commit) Sanitize() Commit {
	c.Hash = cleanString(c.Hash)
	c.Type = cleanString(c.Type)
	c.Component = cleanString(c.Component)
	c.Remote = cleanString(c.Remote)
	c.Repository = cleanString(c.Repository)

	return c
}

// Check verifies that Commit is valid
func (c Commit) Check() error {
	if len(c.Hash) == 0 {
		return fmt.Errorf("commit's hash is required (e.g. `1ab2c3f4d`)")
	}

	if len(c.Type) == 0 {
		return fmt.Errorf("commit's type is required (e.g. `feat`)")
	}

	if len(c.Content) == 0 {
		return fmt.Errorf("commit's content is required (e.g. `Add README.md`)")
	}

	if c.Date == 0 {
		return fmt.Errorf("commit's date is required (e.g. `1596913344`)")
	}

	if len(c.Remote) == 0 {
		return fmt.Errorf("repository's remote is required (e.g. `github.com`)")
	}

	if len(c.Repository) == 0 {
		return fmt.Errorf("repository's name is required (e.g. `vibioh/herodote`)")
	}

	return nil
}

func cleanString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
