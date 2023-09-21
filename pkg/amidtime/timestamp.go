package amidtime

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	raw := strings.Trim(string(data), `"`)
	sec, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}
	if sec == 0 {
		*t = Timestamp{}
	} else {
		*t = Timestamp{time.Unix(sec, 0)}
	}
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`0`), nil
	}
	return []byte(fmt.Sprintf(`%d`, t.Unix())), nil
}

func (tm Timestamp) Value() (driver.Value, error) {
	if tm.IsZero() {
		return nil, nil
	}
	return time.Unix(tm.Unix(), 0), nil
}

func (t *Timestamp) Scan(src any) error {
	switch src := src.(type) {
	case nil:
		*t = Timestamp{}
		return nil
	case time.Time:
		*t = Timestamp{Time: src}
		return nil
	}
	return fmt.Errorf("cannot scan %T", src)
}
