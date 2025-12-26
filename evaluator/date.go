package evaluator

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nulang/nulang/object"
)

// Date represents a JavaScript-like Date object
type Date struct {
	Time time.Time
}

func (d *Date) Type() object.ObjectType { return object.DATE_OBJ }
func (d *Date) Inspect() string         { return d.Time.String() }

// initDateConstructor creates the Date constructor
func initDateConstructor() *object.Builtin {
	return &object.Builtin{Name: "Date", Fn: func(args ...object.Object) object.Object {
		var t time.Time

		if len(args) == 0 {
			// new Date() - current time
			t = time.Now()
		} else if len(args) == 1 {
			switch arg := args[0].(type) {
			case *object.Number:
				// Timestamp in milliseconds
				ms := int64(arg.Value)
				t = time.UnixMilli(ms)
			case *object.String:
				// Parse date string
				parsed, err := parseDate(arg.Value)
				if err != nil {
					t = time.Time{}
				} else {
					t = parsed
				}
			default:
				t = time.Now()
			}
		} else {
			// new Date(year, month, day?, hours?, minutes?, seconds?, ms?)
			year := int(args[0].(*object.Number).Value)
			month := time.Month(int(args[1].(*object.Number).Value) + 1) // JS months are 0-indexed
			day := 1
			hour := 0
			min := 0
			sec := 0
			nsec := 0

			if len(args) > 2 {
				day = int(args[2].(*object.Number).Value)
			}
			if len(args) > 3 {
				hour = int(args[3].(*object.Number).Value)
			}
			if len(args) > 4 {
				min = int(args[4].(*object.Number).Value)
			}
			if len(args) > 5 {
				sec = int(args[5].(*object.Number).Value)
			}
			if len(args) > 6 {
				nsec = int(args[6].(*object.Number).Value) * 1000000
			}

			t = time.Date(year, month, day, hour, min, sec, nsec, time.Local)
		}

		return createDateObject(&Date{Time: t})
	}}
}

// parseDate parses various date string formats
func parseDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC1123,
		time.RFC1123Z,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02/2006",
		"02-Jan-2006",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// createDateObject wraps Date in an ObjectMap with methods
func createDateObject(d *Date) *object.ObjectMap {
	obj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Store the date
	obj.Set("_date", d)

	// getTime() - returns milliseconds since epoch
	obj.Set("getTime", &object.Builtin{Name: "getTime", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// getFullYear()
	obj.Set("getFullYear", &object.Builtin{Name: "getFullYear", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Year())}
	}})

	// getMonth() - 0-11
	obj.Set("getMonth", &object.Builtin{Name: "getMonth", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Month() - 1)}
	}})

	// getDate() - day of month
	obj.Set("getDate", &object.Builtin{Name: "getDate", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Day())}
	}})

	// getDay() - day of week (0 = Sunday)
	obj.Set("getDay", &object.Builtin{Name: "getDay", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Weekday())}
	}})

	// getHours()
	obj.Set("getHours", &object.Builtin{Name: "getHours", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Hour())}
	}})

	// getMinutes()
	obj.Set("getMinutes", &object.Builtin{Name: "getMinutes", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Minute())}
	}})

	// getSeconds()
	obj.Set("getSeconds", &object.Builtin{Name: "getSeconds", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Second())}
	}})

	// getMilliseconds()
	obj.Set("getMilliseconds", &object.Builtin{Name: "getMilliseconds", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.Nanosecond() / 1000000)}
	}})

	// getTimezoneOffset() - returns offset in minutes
	obj.Set("getTimezoneOffset", &object.Builtin{Name: "getTimezoneOffset", Fn: func(args ...object.Object) object.Object {
		_, offset := d.Time.Zone()
		return &object.Number{Value: float64(-offset / 60)}
	}})

	// setTime(ms)
	obj.Set("setTime", &object.Builtin{Name: "setTime", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: float64(d.Time.UnixMilli())}
		}
		if num, ok := args[0].(*object.Number); ok {
			d.Time = time.UnixMilli(int64(num.Value))
		}
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// setFullYear(year, month?, day?)
	obj.Set("setFullYear", &object.Builtin{Name: "setFullYear", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: float64(d.Time.UnixMilli())}
		}
		year := int(args[0].(*object.Number).Value)
		month := d.Time.Month()
		day := d.Time.Day()
		if len(args) > 1 {
			month = time.Month(int(args[1].(*object.Number).Value) + 1)
		}
		if len(args) > 2 {
			day = int(args[2].(*object.Number).Value)
		}
		d.Time = time.Date(year, month, day, d.Time.Hour(), d.Time.Minute(), d.Time.Second(), d.Time.Nanosecond(), d.Time.Location())
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// setMonth(month, day?)
	obj.Set("setMonth", &object.Builtin{Name: "setMonth", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: float64(d.Time.UnixMilli())}
		}
		month := time.Month(int(args[0].(*object.Number).Value) + 1)
		day := d.Time.Day()
		if len(args) > 1 {
			day = int(args[1].(*object.Number).Value)
		}
		d.Time = time.Date(d.Time.Year(), month, day, d.Time.Hour(), d.Time.Minute(), d.Time.Second(), d.Time.Nanosecond(), d.Time.Location())
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// setDate(day)
	obj.Set("setDate", &object.Builtin{Name: "setDate", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: float64(d.Time.UnixMilli())}
		}
		day := int(args[0].(*object.Number).Value)
		d.Time = time.Date(d.Time.Year(), d.Time.Month(), day, d.Time.Hour(), d.Time.Minute(), d.Time.Second(), d.Time.Nanosecond(), d.Time.Location())
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// setHours(hours, min?, sec?, ms?)
	obj.Set("setHours", &object.Builtin{Name: "setHours", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: float64(d.Time.UnixMilli())}
		}
		hour := int(args[0].(*object.Number).Value)
		min := d.Time.Minute()
		sec := d.Time.Second()
		nsec := d.Time.Nanosecond()
		if len(args) > 1 {
			min = int(args[1].(*object.Number).Value)
		}
		if len(args) > 2 {
			sec = int(args[2].(*object.Number).Value)
		}
		if len(args) > 3 {
			nsec = int(args[3].(*object.Number).Value) * 1000000
		}
		d.Time = time.Date(d.Time.Year(), d.Time.Month(), d.Time.Day(), hour, min, sec, nsec, d.Time.Location())
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// toISOString()
	obj.Set("toISOString", &object.Builtin{Name: "toISOString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.UTC().Format(time.RFC3339)}
	}})

	// toDateString()
	obj.Set("toDateString", &object.Builtin{Name: "toDateString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.Format("Mon Jan 02 2006")}
	}})

	// toTimeString()
	obj.Set("toTimeString", &object.Builtin{Name: "toTimeString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.Format("15:04:05 MST")}
	}})

	// toLocaleDateString()
	obj.Set("toLocaleDateString", &object.Builtin{Name: "toLocaleDateString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.Format("01/02/2006")}
	}})

	// toLocaleTimeString()
	obj.Set("toLocaleTimeString", &object.Builtin{Name: "toLocaleTimeString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.Format("03:04:05 PM")}
	}})

	// toLocaleString()
	obj.Set("toLocaleString", &object.Builtin{Name: "toLocaleString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.Format("01/02/2006, 03:04:05 PM")}
	}})

	// toString()
	obj.Set("toString", &object.Builtin{Name: "toString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)")}
	}})

	// toUTCString()
	obj.Set("toUTCString", &object.Builtin{Name: "toUTCString", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.UTC().Format(time.RFC1123)}
	}})

	// toJSON()
	obj.Set("toJSON", &object.Builtin{Name: "toJSON", Fn: func(args ...object.Object) object.Object {
		return &object.String{Value: d.Time.UTC().Format(time.RFC3339)}
	}})

	// valueOf()
	obj.Set("valueOf", &object.Builtin{Name: "valueOf", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(d.Time.UnixMilli())}
	}})

	// format(pattern) - formats date with a pattern (YYYY-MM-DD, etc)
	obj.Set("format", &object.Builtin{Name: "format", Fn: func(args ...object.Object) object.Object {
		pattern := "YYYY-MM-DD HH:mm:ss"
		if len(args) > 0 {
			if str, ok := args[0].(*object.String); ok {
				pattern = str.Value
			}
		}
		return &object.String{Value: formatDate(d, pattern)}
	}})

	// add(duration) - adds duration to date and returns new Date
	obj.Set("add", &object.Builtin{Name: "add", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return obj
		}

		var duration time.Duration
		switch arg := args[0].(type) {
		case *object.Number:
			// Treat as milliseconds
			duration = time.Duration(int64(arg.Value)) * time.Millisecond
		case *object.String:
			// Parse duration string
			parsed, err := parseDuration(arg.Value)
			if err != nil {
				return newError("Invalid duration: %s", arg.Value)
			}
			duration = parsed
		default:
			return obj
		}

		newDate := &Date{Time: d.Time.Add(duration)}
		return createDateObject(newDate)
	}})

	// subtract(duration) - subtracts duration from date and returns new Date
	obj.Set("subtract", &object.Builtin{Name: "subtract", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return obj
		}

		var duration time.Duration
		switch arg := args[0].(type) {
		case *object.Number:
			// Treat as milliseconds
			duration = time.Duration(int64(arg.Value)) * time.Millisecond
		case *object.String:
			// Parse duration string
			parsed, err := parseDuration(arg.Value)
			if err != nil {
				return newError("Invalid duration: %s", arg.Value)
			}
			duration = parsed
		default:
			return obj
		}

		newDate := &Date{Time: d.Time.Add(-duration)}
		return createDateObject(newDate)
	}})

	// diff(otherDate) - returns difference in milliseconds
	obj.Set("diff", &object.Builtin{Name: "diff", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: 0}
		}

		var otherTime time.Time
		switch arg := args[0].(type) {
		case *object.ObjectMap:
			if dateObj, ok := arg.Get("_date"); ok {
				if otherDate, ok := dateObj.(*Date); ok {
					otherTime = otherDate.Time
				}
			}
		case *object.Number:
			// Treat as timestamp
			otherTime = time.UnixMilli(int64(arg.Value))
		default:
			return &object.Number{Value: 0}
		}

		diff := d.Time.Sub(otherTime)
		return &object.Number{Value: float64(diff.Milliseconds())}
	}})

	return obj
}

// initDateStaticMethods adds static methods to Date
func initDateStaticMethods() *object.ObjectMap {
	dateObj := &object.ObjectMap{Pairs: make(map[string]object.ObjectPair)}

	// Date.now() - returns current timestamp
	dateObj.Set("now", &object.Builtin{Name: "now", Fn: func(args ...object.Object) object.Object {
		return &object.Number{Value: float64(time.Now().UnixMilli())}
	}})

	// Date.parse(string) - returns timestamp
	dateObj.Set("parse", &object.Builtin{Name: "parse", Fn: func(args ...object.Object) object.Object {
		if len(args) < 1 {
			return &object.Number{Value: float64(0)}
		}
		str := objectToString(args[0])
		t, err := parseDate(str)
		if err != nil {
			return &object.Number{Value: float64(0)} // NaN in JS, we use 0
		}
		return &object.Number{Value: float64(t.UnixMilli())}
	}})

	// Date.UTC(year, month, day?, hours?, min?, sec?, ms?) - returns timestamp
	dateObj.Set("UTC", &object.Builtin{Name: "UTC", Fn: func(args ...object.Object) object.Object {
		if len(args) < 2 {
			return &object.Number{Value: float64(0)}
		}
		year := int(args[0].(*object.Number).Value)
		month := time.Month(int(args[1].(*object.Number).Value) + 1)
		day := 1
		hour := 0
		min := 0
		sec := 0
		nsec := 0

		if len(args) > 2 {
			day = int(args[2].(*object.Number).Value)
		}
		if len(args) > 3 {
			hour = int(args[3].(*object.Number).Value)
		}
		if len(args) > 4 {
			min = int(args[4].(*object.Number).Value)
		}
		if len(args) > 5 {
			sec = int(args[5].(*object.Number).Value)
		}
		if len(args) > 6 {
			nsec = int(args[6].(*object.Number).Value) * 1000000
		}

		t := time.Date(year, month, day, hour, min, sec, nsec, time.UTC)
		return &object.Number{Value: float64(t.UnixMilli())}
	}})

	// Make Date callable as constructor
	dateObj.Set("__call__", initDateConstructor())

	return dateObj
}

// formatDate formats a date with a pattern
func formatDate(d *Date, pattern string) string {
	// Convert common JS format patterns to Go format
	replacements := []struct {
		js string
		go_ string
	}{
		{"YYYY", "2006"},
		{"YY", "06"},
		{"MMMM", "January"},
		{"MMM", "Jan"},
		{"MM", "01"},
		{"M", "1"},
		{"DDDD", "Monday"},
		{"DDD", "Mon"},
		{"DD", "02"},
		{"D", "2"},
		{"HH", "15"},
		{"hh", "03"},
		{"mm", "04"},
		{"ss", "05"},
		{"SSS", "000"},
		{"A", "PM"},
		{"a", "pm"},
		{"ZZ", "-0700"},
		{"Z", "Z07:00"},
	}

	result := pattern
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.js, r.go_)
	}

	return d.Time.Format(result)
}

// parseDuration parses duration string like "1h", "30m", "2s"
func parseDuration(s string) (time.Duration, error) {
	// Try standard Go duration format first
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Try parsing as milliseconds
	if ms, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Duration(ms) * time.Millisecond, nil
	}

	return 0, fmt.Errorf("invalid duration: %s", s)
}
