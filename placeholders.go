package main

import (
	"regexp"
	"strings"
	"time"
)

type Placeholders struct {
	Year  string
	Month string
	Day   string
}

func NewPlaceholders(d time.Time) *Placeholders {
	return &Placeholders{
		Year:  d.Format("2006"),
		Month: d.Format("01"),
		Day:   d.Format("02"),
	}
}

func (p *Placeholders) Apply(text string) string {
	result := text

	result = strings.ReplaceAll(result, "{year}", p.Year)
	result = strings.ReplaceAll(result, "{month}", p.Month)
	result = strings.ReplaceAll(result, "{day}", p.Day)

	re := regexp.MustCompile(`\{[^}]*\}`)
	result = re.ReplaceAllString(result, "")

	return result
}
