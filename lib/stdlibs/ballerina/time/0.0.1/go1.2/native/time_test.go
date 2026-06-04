// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package native

import (
	"testing"
	"time"

	"ballerina-lang-go/decimal"
)

func TestRoundingUnit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		precision int
		want      int64
		desc      string
	}{
		{0, 1_000_000_000, "precision 0 → 1s in ns"},
		{1, 100_000_000, "precision 1 → 100ms"},
		{2, 10_000_000, "precision 2 → 10ms"},
		{3, 1_000_000, "precision 3 → 1ms"},
		{4, 100_000, "precision 4 → 100µs"},
		{5, 10_000, "precision 5 → 10µs"},
		{6, 1_000, "precision 6 → 1µs"},
		{7, 100, "precision 7 → 100ns"},
		{8, 10, "precision 8 → 10ns"},
		{9, 1, "precision 9 → 1ns"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := roundingUnit(tc.precision)
			if got != tc.want {
				t.Errorf("roundingUnit(%d) = %d, want %d", tc.precision, got, tc.want)
			}
		})
	}
}

func TestFormatInstantNanos(t *testing.T) {
	t.Parallel()
	tests := []struct {
		nanos int
		want  string
		desc  string
	}{
		{0, "", "zero nanos → empty string"},
		{1_000_000, ".001", "1ms → .NNN"},
		{100_000_000, ".100", "100ms → .NNN"},
		{999_000_000, ".999", "999ms → .NNN"},
		{1_000, ".000001", "1µs → .NNNNNN"},
		{100_000, ".000100", "100µs → .NNNNNN"},
		{999_000, ".000999", "999µs → .NNNNNN"},
		{1, ".000000001", "1ns → .NNNNNNNNN"},
		{999_999_999, ".999999999", "max nanos → .NNNNNNNNN"},
		{123_456_789, ".123456789", "mixed nanos → .NNNNNNNNN"},
		{123_000_000, ".123", "exact ms → .NNN"},
		{123_456_000, ".123456", "exact µs → .NNNNNN"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := formatInstantNanos(tc.nanos)
			if got != tc.want {
				t.Errorf("formatInstantNanos(%d) = %q, want %q", tc.nanos, got, tc.want)
			}
		})
	}
}

func TestFormatRFC3339Instant(t *testing.T) {
	t.Parallel()
	tests := []struct {
		t    time.Time
		want string
		desc string
	}{
		{time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), "2024-01-15T10:30:00Z", "zero nanos"},
		{time.Date(2024, 1, 15, 10, 30, 0, 500_000_000, time.UTC), "2024-01-15T10:30:00.500Z", "500ms"},
		{time.Date(2024, 12, 31, 23, 59, 59, 999_999_999, time.UTC), "2024-12-31T23:59:59.999999999Z", "end of year"},
		{time.Date(2024, 6, 1, 0, 0, 0, 1_000_000, time.UTC), "2024-06-01T00:00:00.001Z", "1ms"},
		// Non-UTC input is converted to UTC
		{time.Date(2024, 1, 15, 12, 0, 0, 0, time.FixedZone("+05:30", 5*3600+1800)), "2024-01-15T06:30:00Z", "non-UTC converted"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := formatRFC3339Instant(tc.t)
			if got != tc.want {
				t.Errorf("formatRFC3339Instant(%v) = %q, want %q", tc.t, got, tc.want)
			}
		})
	}
}

func TestFormatRFC3339WithOffset(t *testing.T) {
	t.Parallel()
	base := time.Date(2024, 6, 15, 12, 30, 45, 0, time.UTC)
	tests := []struct {
		t          time.Time
		offsetSecs int
		want       string
		desc       string
	}{
		{base, 0, "2024-06-15T12:30:45Z", "zero offset uses Instant format"},
		{base, 5*3600 + 30*60, "2024-06-15T12:30:45+05:30", "positive offset"},
		{base, -(8 * 3600), "2024-06-15T12:30:45-08:00", "negative offset"},
		{base, 3600, "2024-06-15T12:30:45+01:00", "one hour east"},
		{base, -(5*3600 + 45*60), "2024-06-15T12:30:45-05:45", "Nepal-like negative offset"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := formatRFC3339WithOffset(tc.t, tc.offsetSecs)
			if got != tc.want {
				t.Errorf("formatRFC3339WithOffset(%v, %d) = %q, want %q",
					tc.t, tc.offsetSecs, got, tc.want)
			}
		})
	}
}

func TestOffsetToZoneName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		offsetSecs int
		want       string
		desc       string
	}{
		{0, "Z", "UTC is Z"},
		{3600, "+01:00", "one hour east"},
		{-3600, "-01:00", "one hour west"},
		{5*3600 + 30*60, "+05:30", "India +05:30"},
		{-(5*3600 + 30*60), "-05:30", "negative +05:30"},
		{12 * 3600, "+12:00", "UTC+12"},
		{-(12 * 3600), "-12:00", "UTC-12"},
		{30 * 60, "+00:30", "half hour east"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := offsetToZoneName(tc.offsetSecs)
			if got != tc.want {
				t.Errorf("offsetToZoneName(%d) = %q, want %q", tc.offsetSecs, got, tc.want)
			}
		})
	}
}

func TestZoneAbbrevFor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		t    time.Time
		want string
		desc string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Z", "UTC maps to Z"},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.FixedZone("+05:30", 5*3600+1800)), "+05:30", "fixed zone preserved"},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.FixedZone("-08:00", -8*3600)), "-08:00", "negative fixed zone"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := zoneAbbrevFor(tc.t)
			if got != tc.want {
				t.Errorf("zoneAbbrevFor(%v) = %q, want %q", tc.t, got, tc.want)
			}
		})
	}
}

func TestNormaliseEmailDay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
		desc  string
	}{
		{"Mon, 3 Dec 2024 10:00:00 +0000", "Mon, 03 Dec 2024 10:00:00 +0000", "single digit day padded"},
		{"Mon, 03 Dec 2024 10:00:00 +0000", "Mon, 03 Dec 2024 10:00:00 +0000", "two-digit day unchanged"},
		{"Mon, 15 Jan 2024 10:00:00 +0000", "Mon, 15 Jan 2024 10:00:00 +0000", "double digit day unchanged"},
		{"Mon, 1 Jan 2024 10:00:00 +0000", "Mon, 01 Jan 2024 10:00:00 +0000", "pad day 1"},
		{"Mon, 9 Jun 2024 00:00:00 +0000", "Mon, 09 Jun 2024 00:00:00 +0000", "pad day 9"},
		{"", "", "empty string unchanged"},
		{"short", "short", "too short unchanged"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := normaliseEmailDay(tc.input)
			if got != tc.want {
				t.Errorf("normaliseEmailDay(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestEmailOffsetStr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
		desc  string
	}{
		{"0", "+0000", "zero maps to +0000"},
		{"+05:30", "+05:30", "non-zero preserved"},
		{"-08:00", "-08:00", "negative preserved"},
		{"+0000", "+0000", "already +0000"},
		{"UTC", "UTC", "UTC preserved"},
		{"Z", "Z", "Z preserved"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := emailOffsetStr(tc.input)
			if got != tc.want {
				t.Errorf("emailOffsetStr(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsNumericOffset(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  bool
		desc  string
	}{
		{"+05:30", true, "positive offset"},
		{"-08:00", true, "negative offset"},
		{"+00:00", true, "zero offset"},
		{"+12:45", true, "large positive"},
		{"UTC", false, "UTC is not numeric"},
		{"Z", false, "Z too short"},
		{"America/New_York", false, "IANA name is not numeric"},
		{"", false, "empty string"},
		{"ab", false, "no sign"},
		{"+abc", false, "non-digit after sign"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := isNumericOffset(tc.input)
			if got != tc.want {
				t.Errorf("isNumericOffset(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseNumericOffset(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		want    int
		wantErr bool
		desc    string
	}{
		{"+05:30", 5*3600 + 30*60, false, "positive +05:30"},
		{"-08:00", -(8 * 3600), false, "negative -08:00"},
		{"+00:00", 0, false, "zero offset"},
		{"+12:00", 12 * 3600, false, "UTC+12"},
		{"-05:45", -(5*3600 + 45*60), false, "negative with minutes"},
		{"+0530", 5*3600 + 30*60, false, "no-colon format HHMM"},
		{"-0800", -(8 * 3600), false, "no-colon negative"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := parseNumericOffset(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("parseNumericOffset(%q): expected error", tc.input)
				}
				return
			}
			if err != nil {
				t.Errorf("parseNumericOffset(%q): unexpected error: %v", tc.input, err)
				return
			}
			if got != tc.want {
				t.Errorf("parseNumericOffset(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestNanosToFrac(t *testing.T) {
	t.Parallel()
	tests := []struct {
		nanos   int64
		wantF64 float64
		desc    string
	}{
		{0, 0.0, "zero nanos"},
		{500_000_000, 0.5, "half second"},
		{1_000_000, 0.001, "1ms"},
		{1, 1e-9, "1ns"},
		{999_999_999, 0.999999999, "max nanos"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := nanosToFrac(tc.nanos)
			gotF := got.Float64()
			// Allow small floating-point tolerance
			diff := gotF - tc.wantF64
			if diff < 0 {
				diff = -diff
			}
			if diff > 1e-15 {
				t.Errorf("nanosToFrac(%d).Float64() = %v, want %v", tc.nanos, gotF, tc.wantF64)
			}
		})
	}
}

func TestDecimalToSecNano(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    *decimal.Decimal
		wantSec  int
		wantNano int
		desc     string
	}{
		{decimal.FromInt64(0), 0, 0, "zero"},
		{decimal.FromInt64(1), 1, 0, "1 second"},
		{decimal.FromInt64(5), 5, 0, "5 seconds"},
		{nil, 0, 0, "nil input"},
	}
	// Build decimal 1.5 from string for exact representation
	if d, err := decimal.FromString("1.5"); err == nil {
		tests = append(tests, struct {
			input    *decimal.Decimal
			wantSec  int
			wantNano int
			desc     string
		}{d, 1, 500_000_000, "1.5 seconds"})
	}
	if d, err := decimal.FromString("0.001"); err == nil {
		tests = append(tests, struct {
			input    *decimal.Decimal
			wantSec  int
			wantNano int
			desc     string
		}{d, 0, 1_000_000, "1ms"})
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotSec, gotNano := decimalToSecNano(tc.input)
			if gotSec != tc.wantSec {
				t.Errorf("decimalToSecNano(%v) sec = %d, want %d", tc.input, gotSec, tc.wantSec)
			}
			if gotNano != tc.wantNano {
				t.Errorf("decimalToSecNano(%v) nano = %d, want %d", tc.input, gotNano, tc.wantNano)
			}
		})
	}
}

func TestFormatEmailDate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		t         time.Time
		offsetStr string
		want      string
		desc      string
	}{
		{
			time.Date(2024, time.December, 5, 10, 30, 45, 0, time.UTC),
			"+0000",
			"Thu, 5 Dec 2024 10:30:45 +0000",
			"basic date",
		},
		{
			time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			"+0530",
			"Mon, 1 Jan 2024 00:00:00 +0530",
			"new year",
		},
		{
			time.Date(2024, time.June, 15, 23, 59, 59, 0, time.UTC),
			"-0800",
			"Sat, 15 Jun 2024 23:59:59 -0800",
			"negative offset",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := formatEmailDate(tc.t, tc.offsetStr)
			if got != tc.want {
				t.Errorf("formatEmailDate(%v, %q) = %q, want %q",
					tc.t, tc.offsetStr, got, tc.want)
			}
		})
	}
}

func TestParseEmailDate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		wantErr bool
		desc    string
	}{
		{"Thu, 05 Dec 2024 10:30:45 +0000", false, "two-digit day"},
		{"Mon, 01 Jan 2024 00:00:00 +0530", false, "positive offset"},
		{"Sat, 15 Jun 2024 23:59:59 -0800", false, "negative offset"},
		{"not a date", true, "invalid input"},
		{"", true, "empty string"},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := parseEmailDate(tc.input)
			if tc.wantErr && err == nil {
				t.Errorf("parseEmailDate(%q): expected error, got nil", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("parseEmailDate(%q): unexpected error: %v", tc.input, err)
			}
		})
	}
}

func TestParseEmailDateRoundtrip(t *testing.T) {
	t.Parallel()
	// Format a known time and re-parse it; verify the round-trip.
	original := time.Date(2024, time.June, 15, 10, 30, 45, 0, time.UTC)
	formatted := formatEmailDate(original, "+0000")
	// normalise to 2-digit day for parseEmailDate
	normalised := normaliseEmailDay(formatted)
	parsed, err := parseEmailDate(normalised)
	if err != nil {
		t.Fatalf("parseEmailDate(%q): %v", normalised, err)
	}
	if !parsed.Equal(original) {
		t.Errorf("roundtrip: got %v, want %v", parsed, original)
	}
}
