/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package iec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

// udt mimics a user-defined type holding every time type, as sent by honeycomb.
type udt struct {
	DT   DT
	DATE DATE
	TOD  TOD
	TS   TIMESPEC
	TIME TIME
}

var when = time.Date(2026, 10, 1, 14, 30, 5, 0, time.UTC)

func todOf(h, m, s, ns int) TOD {
	return TOD(time.Time{}.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute +
		time.Duration(s)*time.Second + time.Duration(ns)))
}

func TestTimeTypesJSONEncoding(t *testing.T) {
	v := udt{
		DT:   DT(when),
		DATE: DATE(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)),
		TOD:  todOf(14, 30, 5, 0),
		TS:   TIMESPEC(when.Add(500 * time.Millisecond)),
		TIME: TIME(1500 * time.Millisecond),
	}
	got, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"DT":"2026-10-01T14:30:05Z","DATE":"2026-10-01","TOD":"14:30:05","TS":"2026-10-01T14:30:05.5Z","TIME":1500000000}`
	if string(got) != want {
		t.Fatalf("Marshal =\n %s\nwant\n %s", got, want)
	}

	var back udt
	if err := json.Unmarshal(got, &back); err != nil {
		t.Fatal(err)
	}
	if !time.Time(back.DT).Equal(time.Time(v.DT)) || !time.Time(back.DATE).Equal(time.Time(v.DATE)) ||
		!time.Time(back.TOD).Equal(time.Time(v.TOD)) || !time.Time(back.TS).Equal(time.Time(v.TS)) || back.TIME != v.TIME {
		t.Fatalf("round trip = %+v, want %+v", back, v)
	}
}

func TestTimeTypesPreserveZoneAndFraction(t *testing.T) {
	zone := time.FixedZone("CEST", 2*3600)
	in := DT(time.Date(2026, 10, 1, 16, 30, 5, 123456789, zone))
	got, _ := json.Marshal(in)
	if string(got) != `"2026-10-01T16:30:05.123456789+02:00"` {
		t.Fatalf("Marshal = %s", got)
	}
	tod, _ := json.Marshal(todOf(7, 5, 3, 250_000_000))
	if string(tod) != `"07:05:03.25"` {
		t.Fatalf("TOD Marshal = %s", tod)
	}
}

func TestTimeTypesDecodeIECLiterals(t *testing.T) {
	cases := []struct {
		json string
		dst  any
		want time.Time
	}{
		{`"DT#2026-10-01-14:30:05"`, new(DT), when},
		{`"date_and_time#2026-10-01-14:30:05.5"`, new(DT), when.Add(500 * time.Millisecond)},
		{`"2026-10-01T14:30:05"`, new(DT), when}, // no zone: UTC
		{`"D#2026-10-01"`, new(DATE), time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
		{`"DATE#2026-10-01"`, new(DATE), time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
		{`"TOD#14:30:05.5"`, new(TOD), time.Time(todOf(14, 30, 5, 500_000_000))},
		{`"TIME_OF_DAY#14:30:05"`, new(TOD), time.Time(todOf(14, 30, 5, 0))},
		{`"DT#2026-10-01-14:30:05"`, new(TIMESPEC), when},
	}
	for _, c := range cases {
		if err := json.Unmarshal([]byte(c.json), c.dst); err != nil {
			t.Errorf("%s: %v", c.json, err)
			continue
		}
		var got time.Time
		switch d := c.dst.(type) {
		case *DT:
			got = time.Time(*d)
		case *DATE:
			got = time.Time(*d)
		case *TOD:
			got = time.Time(*d)
		case *TIMESPEC:
			got = time.Time(*d)
		}
		if !got.Equal(c.want) {
			t.Errorf("%s decoded to %v, want %v", c.json, got, c.want)
		}
	}
}

func TestDecodedTODMatchesInitDate(t *testing.T) {
	var tod TOD
	if err := json.Unmarshal([]byte(`"00:00:00"`), &tod); err != nil {
		t.Fatal(err)
	}
	if !time.Time(tod).Equal(time.Time(INITTOD)) {
		t.Fatalf("midnight TOD = %v, want INITTOD %v", time.Time(tod), time.Time(INITTOD))
	}
}

func TestTimeTypesRejectBadInput(t *testing.T) {
	cases := []struct {
		json string
		dst  any
	}{
		{`"yesterday"`, new(DT)},
		{`"DT#2026-10-01"`, new(DT)},
		{`"2026-13-01"`, new(DATE)},
		{`"25:00:00"`, new(TOD)},
		{`12`, new(DATE)},
	}
	for _, c := range cases {
		if err := json.Unmarshal([]byte(c.json), c.dst); err == nil {
			t.Errorf("%s into %T: expected an error", c.json, c.dst)
		}
	}
	err := json.Unmarshal([]byte(`"bad"`), new(DT))
	if err == nil || !strings.Contains(err.Error(), "invalid DT") {
		t.Errorf("error = %v, want it to name the type", err)
	}
}

func TestTimeTypesNullLeavesValue(t *testing.T) {
	v := udt{DT: DT(when)}
	if err := json.Unmarshal([]byte(`{"DT":null}`), &v); err != nil {
		t.Fatal(err)
	}
	if !time.Time(v.DT).Equal(when) {
		t.Fatalf("null changed the value to %v", time.Time(v.DT))
	}
}

// TestLegacyEmptyObjectStillDecodes covers data written by v0.1.0-beta1 and
// earlier, which encoded these types as {}. Such rows must keep loading.
func TestLegacyEmptyObjectStillDecodes(t *testing.T) {
	v := udt{DT: DT(when), TIME: 7}
	legacy := `{"DT":{},"DATE":{ },"TOD":{},"TS":{},"TIME":1500000000}`
	if err := json.Unmarshal([]byte(legacy), &v); err != nil {
		t.Fatalf("legacy row failed to decode: %v", err)
	}
	if !time.Time(v.DT).Equal(when) || v.TIME != TIME(1500*time.Millisecond) {
		t.Fatalf("decoded %+v: legacy {} must leave the value unchanged and other fields must decode", v)
	}
	// Any other object or a non-string value is still an error.
	for _, bad := range []string{`{"DT":{"wall":1}}`, `{"DT":12}`, `{"DT":true}`} {
		if err := json.Unmarshal([]byte(bad), new(udt)); err == nil {
			t.Errorf("%s: expected an error", bad)
		}
	}
}

func TestTimeTypesGobRoundTrip(t *testing.T) {
	zone := time.FixedZone("EST", -5*3600)
	in := udt{
		DT:   DT(when.In(zone)),
		DATE: DATE(when),
		TOD:  todOf(1, 2, 3, 4),
		TS:   TIMESPEC(when),
		TIME: TIME(time.Second),
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		t.Fatal(err)
	}
	var out udt
	if err := gob.NewDecoder(&buf).Decode(&out); err != nil {
		t.Fatal(err)
	}
	_, inOff := time.Time(in.DT).Zone()
	_, outOff := time.Time(out.DT).Zone()
	if !time.Time(out.DT).Equal(time.Time(in.DT)) || inOff != outOff ||
		!time.Time(out.TOD).Equal(time.Time(in.TOD)) || out.TIME != in.TIME {
		t.Fatalf("gob round trip = %+v, want %+v", out, in)
	}
}

func ExampleDT_MarshalText() {
	type Reading struct {
		Value REAL
		At    DT
	}
	b, _ := json.Marshal(Reading{Value: 21.5, At: DT(time.Date(2026, 10, 1, 14, 30, 5, 0, time.UTC))})
	fmt.Println(string(b))

	var r Reading
	_ = json.Unmarshal([]byte(`{"Value":3,"At":"DT#2026-10-01-08:00:00"}`), &r)
	fmt.Println(time.Time(r.At).Format(time.RFC3339))
	// Output:
	// {"Value":21.5,"At":"2026-10-01T14:30:05Z"}
	// 2026-10-01T08:00:00Z
}
