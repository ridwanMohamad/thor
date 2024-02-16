package date_util

import (
	"gopkg.in/guregu/null.v4"
	"thor/src/constanta"
	"time"
)

func StringToNilTime(tm string) (resp null.Time) {
	if !null.StringFrom(tm).Valid || tm == "" {
		return null.TimeFromPtr(nil)
	}

	if dt, err := time.Parse("2006-01-02", tm); err == nil {
		return null.TimeFrom(dt)
	}

	return null.TimeFromPtr(nil)
}

func T1BelowT2(tm1 string, tm2 string) (rtm1 null.Time, rtm2 null.Time, err error) {
	var dt1 time.Time
	var dt2 time.Time

	if dt1, err = time.Parse("2006-01-02", tm1); err != nil {
		return null.TimeFromPtr(nil), null.TimeFromPtr(nil), err
	}

	if dt2, err = time.Parse("2006-01-02", tm2); err != nil {
		return null.TimeFromPtr(nil), null.TimeFromPtr(nil), err
	}

	if dt2.Before(dt1) == true {
		return null.TimeFromPtr(nil), null.TimeFromPtr(nil), constanta.Date2BeforeDate1
	}

	if dt1.After(dt2) == true {
		return null.TimeFromPtr(nil), null.TimeFromPtr(nil), constanta.Date1AfterDate2
	}
	return null.TimeFrom(dt1), null.TimeFrom(dt2), nil
}
