/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software

 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 */

package royaljelly

import "time"

func (t *TIMESPEC) AFTER(u TIMESPEC) BOOL {
	return BOOL(time.Time(*t).After(time.Time(u)))
}

func (t *TIMESPEC) BEFORE(u TIMESPEC) BOOL {
	return BOOL(time.Time(*t).Before(time.Time(u)))
}

func (t *TIMESPEC) EQUAL(u TIMESPEC) BOOL {
	return BOOL(time.Time(*t).Equal(time.Time(u)))
}

func (t *TIMESPEC) MONTH_STRING() STRING {
	return STRING(time.Time(*t).Month().String())
}

func (t *TIMESPEC) WEEKDAY_STRING() STRING {
	return STRING(time.Time(*t).Weekday().String())
}

func (t *TIMESPEC) ISZERO() BOOL {
	return BOOL(time.Time(*t).IsZero())
}

func (t *TIMESPEC) TIME() (tm TIME) {
	// This conversion is conceptually tricky. IEC TIME is a duration.
	// A TIMESPEC is a point in time. A common interpretation is the duration since midnight.
	t_time := time.Time(*t)
	midnight := time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location())
	return TIME(t_time.Sub(midnight))
}

func (t *TIMESPEC) DATE() (d DATE) {
	t_time := time.Time(*t)
	return DATE(time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location()))
}

func (t *TIMESPEC) CLOCK() (tod TOD) {
	return TOD(time.Time(*t))
}

func (t *TIMESPEC) DATETIME() (dt DT) {
	return DT(time.Time(*t))
}

func (t *TIMESPEC) YEAR() ANYINT {
	return ANYINT(time.Time(*t).Year())
}

func (t *TIMESPEC) MONTH() ANYINT {
	return ANYINT(time.Time(*t).Month())
}

func (t *TIMESPEC) HOUR() ANYINT {
	return ANYINT(time.Time(*t).Hour())
}

func (t *TIMESPEC) MINUTE() ANYINT {
	return ANYINT(time.Time(*t).Minute())
}

func (t *TIMESPEC) MILLISECOND() ANYINT {
	return ANYINT(time.Time(*t).Nanosecond() / 1e6)
}

func (t *TIMESPEC) SECOND() ANYINT {
	return ANYINT(time.Time(*t).Second())
}

func (t *TIMESPEC) DAY() ANYINT {
	return ANYINT(time.Time(*t).Day())
}

func (t *TIMESPEC) ISOWEEK() (year, week ANYINT) {
	y, w := time.Time(*t).ISOWeek()
	return ANYINT(y), ANYINT(w)
}

func (t *TIMESPEC) WEEKDAY() ANYINT {
	return ANYINT(time.Time(*t).Weekday())
}

func NOW() TIMESPEC {
	return TIMESPEC(time.Now())
}

func UTC(t TIMESPEC) TIMESPEC {
	return TIMESPEC(time.Time(t).UTC())
}

func LOCAL(t TIMESPEC) TIMESPEC {
	return TIMESPEC(time.Time(t).Local())
}

/*
TO_DT and other conversions
*/

// DT_TO_TM extracts the components of a DT into a TM struct.
func DT_TO_TM(in DT) TM {
	t := time.Time(in)
	return TM{
		d:  t.Day(),
		h:  t.Hour(),
		m:  t.Minute(),
		s:  t.Second(),
		ms: t.Nanosecond() / 1e6,
	}
}

// TM_TO_DT converts a TM struct into a DT (DATE_AND_TIME).
// It uses the current year and month, which is a common approach when only time components are provided.
func TM_TO_DT(in TM) DT {
	now := time.Now()
	return DT(time.Date(now.Year(), now.Month(), in.d, in.h, in.m, in.s, in.ms*1e6, now.Location()))
}

func TOD_TO_DT(in TOD) DT {
	return DT(time.Time(in))
}

func DATE_TO_DT(in DATE) DT {
	return DT(time.Time(in))
}
