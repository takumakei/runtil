package app

import (
	"time"
)

// HHMMSS represents a duration from 00:00:00 in UTC.
type HHMMSS time.Duration

// Parse parses a string s into hms.
func (hms *HHMMSS) Parse(s string) (err error) {
	ll := []string{
		"15:04:05.999999999-07:00",
		"15:04:05.999999999Z07:00",
		"15:04:05.999999999",
		"15:04:05.999999-07:00",
		"15:04:05.999999Z07:00",
		"15:04:05.999999",
		"15:04:05.999-07:00",
		"15:04:05.999Z07:00",
		"15:04:05.999",
		"15:04:05-07:00",
		"15:04:05Z07:00",
		"15:04:05",
		"15:04-07:00",
		"15:04Z07:00",
		"15:04",
		"15-07:00",
		"15Z07:00",
		"15",
	}
	const year = "2006 "
	ys := year + s
	for _, l := range ll {
		tm, e := time.ParseInLocation(year+l, ys, time.Local)
		if e == nil {
			*(*time.Duration)(hms) = tm.Sub(tm.Truncate(24 * time.Hour))
			return nil
		}
		err = e
	}
	err = &time.ParseError{
		Value:   s,
		Message: ": cannot parse into HHMMSS",
	}
	return
}

// String returns a string representation of hms.
func (hms HHMMSS) String() string {
	return time.Duration(hms).String()
}

// Next returns the time.Time that meets hms and comes first after t.
func (hms HHMMSS) Next(t time.Time) time.Time {
	a := t.In(time.UTC)
	b := a.Truncate(24 * time.Hour)
	c := time.Duration(hms)
	if d := a.Sub(b); c <= d {
		c += 24 * time.Hour
	}
	return b.Add(c)
}

// Truncate returns the result of rounding t down to a multiple of d
// (since 00:00 in tm.Location())
func Truncate(tm time.Time, d time.Duration) time.Time {
	offset := func(tm time.Time) time.Duration {
		_, offset := tm.Zone()
		return time.Duration(offset) * time.Second
	}(tm)
	return tm.Add(offset).Truncate(d).Add(-offset)
}
