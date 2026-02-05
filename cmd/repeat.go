package cmd

import (
	"fmt"
	"strings"
	"time"
)

type repeatUnit int

const (
	repeatNone repeatUnit = iota
	repeatDaily
	repeatWeekly
	repeatMonthly
	repeatYearly
)

// parseRepeatUnit parses a repeat pattern string into a repeatUnit.
// Accepts: daily, day, weekly, week, monthly, month, yearly, year, annually
func parseRepeatUnit(raw string) (repeatUnit, error) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return repeatNone, nil
	}
	switch raw {
	case "daily", "day":
		return repeatDaily, nil
	case "weekly", "week":
		return repeatWeekly, nil
	case "monthly", "month":
		return repeatMonthly, nil
	case "yearly", "year", "annually":
		return repeatYearly, nil
	default:
		return repeatNone, fmt.Errorf("invalid repeat value %q (must be daily, weekly, monthly, or yearly)", raw)
	}
}

// expandRepeatSchedule generates a list of dates based on the repeat pattern.
// If both count and until are provided, stops at whichever limit is reached first.
// If neither is provided, returns a single occurrence (the start date).
func expandRepeatSchedule(start time.Time, unit repeatUnit, count int, until *time.Time) []time.Time {
	if unit == repeatNone {
		return []time.Time{start}
	}
	if count < 0 {
		count = 0
	}
	// Defensive guard: if neither count nor until is set, return single occurrence
	if count == 0 && until == nil {
		return []time.Time{start}
	}

	out := []time.Time{}
	for i := 0; ; i++ {
		t := addRepeat(start, unit, i)
		if until != nil && t.After(*until) {
			break
		}
		out = append(out, t)
		if count > 0 && len(out) >= count {
			break
		}
	}
	return out
}

// addRepeat adds n units of the repeat pattern to the given time.
func addRepeat(t time.Time, unit repeatUnit, n int) time.Time {
	switch unit {
	case repeatDaily:
		return t.AddDate(0, 0, n)
	case repeatWeekly:
		return t.AddDate(0, 0, 7*n)
	case repeatMonthly:
		return t.AddDate(0, n, 0)
	case repeatYearly:
		return t.AddDate(n, 0, 0)
	default:
		return t
	}
}
