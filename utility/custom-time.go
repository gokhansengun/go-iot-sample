package utility

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime TODO requires a comment
type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05.000Z"

// UnmarshalJSON TODO requires a comment
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}

// MarshalJSON TODO requires a comment
func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

var nilTime = (time.Time{}).UnixNano()
