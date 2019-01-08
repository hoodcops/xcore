package db

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// NullableTime represent a datetime column
// that can be null
type NullableTime struct {
	mysql.NullTime
}

// MarshalJSON determines how a NullableTime is
// marshalled into JSON
func (nt *NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}
