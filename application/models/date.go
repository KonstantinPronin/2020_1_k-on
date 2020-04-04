package models

import (
	"encoding/json"
	"strings"
	"time"
)

//Type Date with specific time format
type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	t := time.Time(*d).Format("2006-01-02")
	return json.Marshal(t)
}

func (d *Date) IsZero() bool {
	t := time.Time(*d)
	return t.IsZero()
}

func (d *Date) String() string {
	return time.Time(*d).Format("2006-01-02")
}
