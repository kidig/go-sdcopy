package main

import (
	"testing"
	"time"
)

func TestNewPlaceholders(t *testing.T) {
	d := time.Date(2023, time.March, 14, 0, 0, 0, 0, time.UTC)
	p := NewPlaceholders(d)

	if p.Year != "2023" {
		t.Errorf("expected Year to be 2023, got %s", p.Year)
	}
	if p.Month != "03" {
		t.Errorf("expected Month to be 03, got %s", p.Month)
	}
	if p.Day != "14" {
		t.Errorf("expected Day to be 14, got %s", p.Day)
	}
}

func TestApply(t *testing.T) {
	d := time.Date(2023, time.March, 14, 0, 0, 0, 0, time.UTC)
	p := NewPlaceholders(d)

	tests := []struct {
		input    string
		expected string
	}{
		{"Today is {year}-{month}-{day}", "Today is 2023-03-14"},
		{"Year: {year}, Month: {month}, Day: {day}", "Year: 2023, Month: 03, Day: 14"},
		{"No placeholders here", "No placeholders here"},
		{"Unknown {placeholder}", "Unknown "},
	}

	for _, test := range tests {
		result := p.Apply(test.input)
		if result != test.expected {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
	}
}
