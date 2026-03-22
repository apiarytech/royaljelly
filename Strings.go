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

import (
	"fmt"
	"strings"
)

/*********************************/
/* IEC 61131-3 Standard Functions*/
/*********************************/

// LEN returns the length of the input string.
// Conforms to IEC 61131-3, Table 29.
func LEN(s STRING) LINT {
	return LINT(len(string(s)))
}

// LEFT returns the leftmost L characters of the input string IN.
// Conforms to IEC 61131-3, Table 29.
func LEFT(IN STRING, L LINT) STRING {
	s := string(IN)
	if L < 0 {
		panic("LEFT: length L cannot be negative")
	}
	if L >= LINT(len(s)) {
		return IN // Return the whole string if L is greater than or equal to length
	}
	return STRING(s[:L])
}

// RIGHT returns the rightmost L characters of the input string IN.
// Conforms to IEC 61131-3, Table 29.
func RIGHT(IN STRING, L LINT) STRING {
	s := string(IN)
	if L < 0 {
		panic("RIGHT: length L cannot be negative")
	}
	if L >= LINT(len(s)) {
		return IN // Return the whole string if L is greater than or equal to length
	}
	return STRING(s[len(s)-int(L):])
}

// MID returns L characters of the input string IN, beginning at the P-th character.
// P is 1-based, as per the IEC standard.
// Conforms to IEC 61131-3, Table 29.
func MID(IN STRING, L, P LINT) STRING {
	s := string(IN)
	p_zero_based := int(P - 1)
	l_int := int(L)

	if L < 0 {
		panic("MID: length L cannot be negative")
	}
	if P < 1 {
		panic("MID: position P cannot be less than 1")
	}
	if p_zero_based >= len(s) {
		return "" // Position is out of bounds
	}
	if p_zero_based+l_int > len(s) {
		return STRING(s[p_zero_based:]) // Return from P to the end of the string
	}
	return STRING(s[p_zero_based : p_zero_based+l_int])
}

// CONCAT performs extensible concatenation of two or more strings.
// Conforms to IEC 61131-3, Table 29.
func CONCAT(inputs ...STRING) STRING {
	var builder strings.Builder
	for _, s := range inputs {
		builder.WriteString(string(s))
	}
	return STRING(builder.String())
}

// INSERT inserts string IN2 into string IN1 after the P-th character position.
// P is 1-based, as per the IEC standard.
// Conforms to IEC 61131-3, Table 29.
func INSERT(IN1, IN2 STRING, P LINT) STRING {
	s1 := string(IN1)
	p_zero_based := int(P) // Insert *after* P, so index is P
	if P < 0 {
		panic("INSERT: position P cannot be negative")
	}
	if p_zero_based > len(s1) {
		p_zero_based = len(s1) // If P is out of bounds, append to the end
	}
	return STRING(s1[:p_zero_based] + string(IN2) + s1[p_zero_based:])
}

// DELETE deletes L characters from string IN, beginning at the P-th character.
// P is 1-based, as per the IEC standard.
// Conforms to IEC 61131-3, Table 29.
func DELETE(IN STRING, L, P LINT) STRING {
	s := string(IN)
	p_zero_based := int(P - 1)
	l_int := int(L)

	if L < 0 || P < 1 || p_zero_based >= len(s) {
		return IN // Return original string on invalid input
	}
	if p_zero_based+l_int > len(s) {
		return STRING(s[:p_zero_based]) // Delete from P to the end
	}
	return STRING(s[:p_zero_based] + s[p_zero_based+l_int:])
}

// FIND finds the character position of the beginning of the first occurrence of IN2 in IN1.
// Returns a 1-based index, or 0 if not found, as per the IEC standard.
// Conforms to IEC 61131-3, Table 29.
func FIND(IN1, IN2 STRING) LINT {
	index := strings.Index(string(IN1), string(IN2))
	if index == -1 {
		return 0 // Not found
	}
	return LINT(index + 1) // Convert 0-based to 1-based
}

/*********************************/
/*  IEC Strings definitions*/
/*********************************/

// CLONE returns a fresh copy of s. It guarantees to make a copy of s into a new allocation, which can be important when retaining only a small substring of a much larger string. Using Clone can help such programs use less memory. Of course, since using Clone makes a copy, overuse of Clone can make programs use more memory. Clone should typically be used only rarely, and only when profiling indicates that it is needed. For strings of length zero the string "" will be returned and no allocation is made.
func CLONE(s STRING) STRING {
	return STRING(strings.Clone(string(s)))
}

// COMPARE returns an integer comparing two strings lexicographically. The result will be 0 if a == b, -1 if a < b, and +1 if a > b. COMPARE is included only for symmetry with package bytes. It is usually clearer and always faster to use the built-in string comparison operators ==, <, >, and so on.
func COMPARE(a, b STRING) LINT {
	return LINT(strings.Compare(string(a), string(b)))
}

// CONTAINS reports whether substr is within s.
func CONTAINS(s, substr STRING) BOOL {
	return BOOL(strings.Contains(string(s), string(substr)))
}

// CONTAINSANY reports whether any Unicode code points in chars are within s.
func CONTAINSANY(s, chars STRING) BOOL {
	return BOOL(strings.ContainsAny(string(s), string(chars)))
}

// CONTAINSRUNE reports whether the Unicode code point r is within s.
func CONTAINSRUNE(s STRING, r WSTRING) BOOL {
	return BOOL(strings.ContainsRune(string(s), rune(r)))
}

// COUNT counts the number of non-overlapping instances of substr in s. If substr is an empty string, Count returns 1 + the number of Unicode code points in s.
func COUNT(s, substr STRING) LINT {
	return LINT(strings.Count(string(s), string(substr)))
}

// CUT slices s around the first instance of sep, returning the text before and after sep. The found result reports whether sep appears in s. If sep does not appear in s, cut returns s, "", false.
func CUT(s, sep STRING) (before, after STRING, found BOOL) {
	b, a, f := strings.Cut(string(s), string(sep))
	return STRING(b), STRING(a), BOOL(f)
}

// EQUALFOLD reports whether s and t, interpreted as UTF-8 strings, are equal under simple Unicode case-folding, which is a more general form of case-insensitivity.
func EQUALFOLD(s, t STRING) BOOL {
	return BOOL(strings.EqualFold(string(s), string(t)))
}

// FIELDS splits the string s around each instance of one or more consecutive white space characters, as defined by unicode.IsSpace, returning a slice of substrings of s or an empty slice if s contains only white space.
func FIELDS(s STRING) (out STRINGS) {
	return strings.Fields(string(s))
}

// FIELDSFUNC splits the string s at each run of Unicode code points c satisfying f(c) and returns an array of slices of s. If all code points in s satisfy f(c) or the string is empty, an empty slice is returned. FIELDSFUNC makes no guarantees about the order in which it calls f(c) and assumes that f always returns the same value for a given c.
func FIELDSFUNC(s STRING, f func(c rune) bool) STRINGS {
	return strings.FieldsFunc(string(s), f)
}

// HASPREFIX tests whether the string s begins with prefix.
func HASPREFIX(s, prefix STRING) BOOL {
	return BOOL(strings.HasPrefix(string(s), string(prefix)))
}

// HASSUFFIX tests whether the string s ends with suffix.
func HASSUFFIX(s, suffix STRING) BOOL {
	return BOOL(strings.HasSuffix(string(s), string(suffix)))
}

// INDEX returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func INDEX(s, substr STRING) LINT {
	return LINT(strings.Index(string(s), string(substr)))
}

// INDEXANY returns the index of the first instance of any Unicode code point from chars in s, or -1 if no Unicode code point from chars is present in s.
func INDEXANY(s, chars STRING) LINT {
	return LINT(strings.IndexAny(string(s), string(chars)))
}

// INDEXBYTE returns the index of the first instance of c in s, or -1 if c is not present in s.
func INDEXBYTE(s STRING, c BYTE) LINT {
	return LINT(strings.IndexByte(string(s), c.Value()))
}

// INDEXFUNC returns the index into s of the first Unicode code point satisfying f(c), or -1 if none do.
func INDEXFUNC(s STRING, f func(c rune) bool) LINT {
	return LINT(strings.IndexFunc(string(s), f))
}

// INDEXRUNE returns the index of the first instance of the Unicode code point r, or -1 if rune is not present in s. If r is utf8.RuneError, it returns the first instance of any invalid UTF-8 byte sequence.
func INDEXRUNE(s STRING, r WSTRING) LINT {
	return LINT(strings.IndexRune(string(s), rune(r)))
}

// JOIN concatenates the elements of its first argument to create a single string. The separator string sep is placed between elements in the resulting string.
func JOIN(elems STRINGS, sep STRING) STRING {
	return STRING(strings.Join(elems, string(sep)))
}

// LASTINDEX returns the index of the last instance of substr in s, or -1 if substr is not present in s.
func LASTINDEX(s, substr STRING) LINT {
	return LINT(strings.LastIndex(string(s), string(substr)))
}

// LASTINDEXANY returns the index of the last instance of any Unicode code point from chars in s, or -1 if no Unicode code point from chars is present in s.
func LASTINDEXANY(s, chars STRING) LINT {
	return LINT(strings.LastIndexAny(string(s), string(chars)))
}

// LASTINDIXBYTE returns the index into s of the last Unicode code point satisfying f(c), or -1 if none do.
func LASTINDEXBYTE(s STRING, c BYTE) LINT {
	return LINT(strings.LastIndexByte(string(s), c.Value()))
}

// LASTINDEXFUNC returns the index into s of the last Unicode code point satisfying f(c), or -1 if none do.
func LASTINDEXFUNC(s STRING, f func(rune) bool) LINT {
	return LINT(strings.LastIndexFunc(string(s), f))
}

// MAP returns a copy of the string s with all its characters modified according to the mapping function. If mapping returns a negative value, the character is dropped from the string with no replacement.
func MAP(mapping func(rune) rune, s STRING) STRING {
	return STRING(strings.Map(mapping, string(s)))
}

// REPEAT returns a new string consisting of count copies of the string s.  REPEAT panics if count is negative or if the result of (len(s) * count) overflows.
func REPEAT(s STRING, count DINT) (out STRING) {
	if (count > 0) && ((len(s) * int(count)) < 2147483647) {
		out = STRING(strings.Repeat(string(s), int(count)))
	} else {
		out = ""
	}
	return out
}

// REPLACE replaces L characters of string IN1 by string IN2, starting at the P-th character.
// P is 1-based, as per the IEC standard.
// Conforms to IEC 61131-3, Table 29. Renamed from REPLACE to avoid conflict.
func REPLACE(IN1, IN2 STRING, L, P LINT) STRING {
	s1 := string(IN1)
	l_int := int(L)

	if L < 0 || P < 0 {
		panic(fmt.Sprintf("REPLACE: invalid L (%d) or P (%d) for string of length %d", L, P, len(s1)))
	}

	// As per IEC standard, P=0 is equivalent to P=1 for REPLACE.
	if P == 0 {
		P = 1
	}

	p_zero_based := int(P - 1)

	end_delete := p_zero_based + l_int
	if end_delete > len(s1) {
		end_delete = len(s1)
	}

	// When L=0, it's an insert *after* position P. Otherwise, it's a replace *at* position P.
	if L == 0 {
		// This logic mirrors the INSERT function for P > 0.
		if int(P) > len(s1) {
			return STRING(s1 + string(IN2))
		}
		return STRING(s1[:P] + string(IN2) + s1[P:])
	}

	if p_zero_based > len(s1) {
		return STRING(s1 + string(IN2))
	}

	return STRING(s1[:p_zero_based] + string(IN2) + s1[end_delete:])
}

// REPLACE_STR is a wrapper for Go's strings.Replace. It returns a copy of the string s with the first n non-overlapping instances of old replaced by new. If n < 0, there is no limit on the number of replacements.
func REPLACE_STR(s, old, new STRING, n LINT) STRING {
	return STRING(strings.Replace(string(s), string(old), string(new), int(n)))
}

// REPLACEALL is a wrapper for Go's strings.ReplaceAll. It returns a copy of the string s with all non-overlapping instances of old replaced by new.
func REPLACE_ALL(s, old, new STRING) STRING {
	return STRING(strings.ReplaceAll(string(s), string(old), string(new)))
}

/*
FUNC slices s into all substrings separated by sep and returns a slice of the substrings between those separators.

If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.

If sep is empty, Split splits after each UTF-8 sequence. If both s and sep are empty, Split returns an empty slice.

It is equivalent to SplitN with a count of -1.

To split around the first instance of a separator, see Cut.
*/
func SPLIT(s, sep STRING) STRINGS {
	return strings.Split(string(s), string(sep))
}

/*
SPLITAFTER slices s into all substrings after each instance of sep and returns a slice of those substrings.

If s does not contain sep and sep is not empty, SplitAfter returns a slice of length 1 whose only element is s.

If sep is empty, SplitAfter splits after each UTF-8 sequence. If both s and sep are empty, SplitAfter returns an empty slice.

It is equivalent to SplitAfterN with a count of -1.
*/
func SPLITAFTER(s, sep STRING) STRINGS {
	return strings.SplitAfter(string(s), string(sep))
}

/*
SPLITAFTERN slices s into substrings after each instance of sep and returns a slice of those substrings.

The count determines the number of substrings to return:

n > 0: at most n substrings; the last substring will be the unsplit remainder.
n == 0: the result is nil (zero substrings)
n < 0: all substrings

Edge cases for s and sep (for example, empty strings) are handled as described in the documentation for SplitAfter.
*/
func SPLITAFTERN(s, sep STRING, n LINT) STRINGS {
	return strings.SplitAfterN(string(s), string(sep), int(n))
}

/*
SPLITN slices s into substrings separated by sep and returns a slice of the substrings between those separators.

The count determines the number of substrings to return:

n > 0: at most n substrings; the last substring will be the unsplit remainder.
n == 0: the result is nil (zero substrings)
n < 0: all substrings
Edge cases for s and sep (for example, empty strings) are handled as described in the documentation for Split.

To split around the first instance of a separator, see Cut.
*/

func SPLITN(s, sep STRING, n LINT) STRINGS {
	return strings.SplitN(string(s), string(sep), int(n))
}

// TOLOWER returns s with all Unicode letters mapped to their lower case.
func TOLOWER(s STRING) STRING {
	return STRING(strings.ToLower(string(s)))
}

/* //ToLowerSpecial returns a copy of the string s with all Unicode letters mapped to their lower case using the case mapping specified by c.
func TOLOWERSPECIAL(c unicode.SpecialCase, s STRING) STRING {}
*/

// TOTITLE returns a copy of the string s with all Unicode letters mapped to their Unicode title case.
func TOTITLE(s STRING) STRING {
	return STRING(strings.ToTitle(string(s)))
}

/*
ToTitleSpecial returns a copy of the string s with all Unicode letters mapped to their Unicode title case, giving priority to the special casing rules.
func TOTITLESPECIAL(c unicode.SpecialCase, s STRING) STRING {}
*/

// TOUPPER returns s with all Unicode letters mapped to their upper case.
func TOUPPER(s STRING) STRING {
	return STRING(strings.ToUpper(string(s)))
}

/*
 ToUpperSpecial returns a copy of the string s with all Unicode letters mapped to their upper case using the case mapping specified by c.
func TOUPPERSPECIAL(c unicode.SpecialCase, s STRING) STRING {}
*/

// TOVALIDUTF8 returns a copy of the string s with each run of invalid UTF-8 byte sequences replaced by the replacement string, which may be empty.

func TOVALIDUTF8(s, replacement STRING) STRING {
	return STRING(strings.ToValidUTF8(string(s), string(replacement)))
}

// TRIM returns a slice of the string s with all leading and trailing Unicode code points contained in cutset removed.
func TRIM(s, cutset STRING) STRING {
	return STRING(strings.Trim(string(s), string(cutset)))
}

// TRIMFUNC returns a slice of the string s with all leading and trailing Unicode code points c satisfying f(c) removed.
func TRIMFUNC(s STRING, f func(rune) bool) STRING {
	return STRING(strings.TrimFunc(string(s), f))
}

// TRIMLEFT returns a slice of the string s with all leading Unicode code points contained in cutset removed. To remove a prefix, use TrimPrefix instead.
func TRIMLEFT(s, cutset STRING) STRING {
	return STRING(strings.TrimLeft(string(s), string(cutset)))
}

// TRIMLEFTFUNC returns a slice of the string s with all leading Unicode code points c satisfying f(c) removed.
func TRIMLEFTFUNC(s STRING, f func(rune) bool) STRING {
	return STRING(strings.TrimLeftFunc(string(s), f))
}

// TRIMPREFIX returns s without the provided leading prefix string. If s doesn't start with prefix, s is returned unchanged.
func TRIMPREFIX(s, prefix STRING) STRING {
	return STRING(strings.TrimPrefix(string(s), string(prefix)))
}

// TRIMRIGHT returns a slice of the string s, with all trailing Unicode code points contained in cutset removed.  To remove a suffix, use TrimSuffix instead.
func TRIMRIGHT(s, cutset STRING) STRING {
	return STRING(strings.TrimRight(string(s), string(cutset)))
}

// TRIMRIGHTFUNC returns a slice of the string s with all trailing Unicode code points c satisfying f(c) removed.
func TRIMRIGHTFUNC(s STRING, f func(rune) bool) STRING {
	return STRING(strings.TrimRightFunc(string(s), f))
}

// TRIMSPACE returns a slice of the string s, with all leading and trailing white space removed, as defined by Unicode.
func TRIMSPACE(s STRING) STRING {
	return STRING(strings.TrimSpace(string(s)))
}

// TRIMSUFFIX returns s without the provided trailing suffix string. If s doesn't end with suffix, s is returned unchanged.
func TRIMSUFFIX(s, suffix STRING) STRING {
	return STRING(strings.TrimSuffix(string(s), string(suffix)))
}
