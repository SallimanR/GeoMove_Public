package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

/* From "go-DDD" project */

func timestamptzFromTime(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	ts.Scan(t)
	return ts
}

func timeFromTimestamptz(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

func numericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	// Convert float64 to string first, then scan
	n.Scan(fmt.Sprintf("%.3f", f))
	return n
}

func float64FromNumeric(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}

func StringToPgTime(s string) pgtype.Time {
	if s == "" {
		return pgtype.Time{Valid: false}
	}
	t, err := time.Parse("15:04", s)
	if err != nil {
		return pgtype.Time{Valid: false}
	}
	micros := int64(t.Hour())*3600*1_000_000 +
		int64(t.Minute())*60*1_000_000 +
		int64(t.Second())*1_000_000 +
		int64(t.Nanosecond()/1000)
	return pgtype.Time{Microseconds: micros, Valid: true}
}

// func parseTimeToPgType(timeStr string) (pgtype.Time, error) {
// 	var t time.Time
// 	var err error
//
// 	// Try parsing with different formats
// 	t, err = time.Parse("15:04", timeStr)
// 	if err != nil {
// 		t, err = time.Parse("15:04:05", timeStr)
// 		if err != nil {
// 			return pgtype.Time{}, err
// 		}
// 	}
//
// 	// Calculate microseconds since midnight
// 	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
// 	microseconds := t.Sub(midnight).Microseconds()
//
// 	return pgtype.Time{
// 		Microseconds: microseconds,
// 		Valid:        true,
// 	}, nil
// }
//
// func formatPgTime(t pgtype.Time) string {
// 	if !t.Valid {
// 		return ""
// 	}
// 	// Convert microseconds to time.Duration
// 	duration := time.Duration(t.Microseconds) * time.Microsecond
// 	// Create a base time at midnight
// 	midnight := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
// 	// Add the duration to get the actual time
// 	actualTime := midnight.Add(duration)
// 	// Format as HH:MM
// 	return actualTime.Format("15:04")
// }
