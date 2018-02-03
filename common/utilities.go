package common

import (
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// IndexOf searches a string for the first occurrence of another string, starting
// at a given offset.
// @param str1 String to search.
// @param str2 String to search for.
// @param off Integer offset or -1 if not found.
func IndexOf(str1, str2 string, off int) int {
	index := strings.Index(str1[off:], str2)
	if index == -1 {
		return -1
	}
	return index + off
}

// ContainsWildcards searches string for special characters * and ?
// @param string String to be searched
// @return true if string contains wildcard, false otherwise
func ContainsWildcards(str string) bool {
	prev := ' '
	for _, s := range str {
		if s == '*' || s == '?' {
			if prev != '\\' {
				return true
			}
		}
		prev = s
	}
	return false
}

// ContainsQuestions searches string for special characters ?
// @param string String to be searched
// @return true if string contains wildcard, false otherwise
func ContainsQuestions(str string) bool {
	prev := ' '
	for _, s := range str {
		if s == '?' {
			if prev != '\\' {
				return true
			}
		}
		prev = s
	}
	return false
}

// CountEscapeCharacters counts the number of escape characters in the string beginning and ending
// at the given indices
// @param str string to search
// @param start beginning index
// @param end ending index
// @return number of escape characters in string
func CountEscapeCharacters(str string) (result int) {
	active := false
	for _, s := range str {
		if !active && s == '\\' {
			result++
			active = true
		} else {
			active = false
		}
	}
	return result
}

// LengthWithEscapeCharacters counts the number of characters with escape characters
func LengthWithEscapeCharacters(str string) (result int) {
	return len(str) - CountEscapeCharacters(str)
}

// GetUnescapedColonIndex searches a string for the first unescaped colon and returns the index of that colon
// @param str string to search
// @return index of first unescaped colon, or 0 if not found
func GetUnescapedColonIndex(str string) (idx int) {
	for i, s := range str {
		if s == ':' {
			if i == 0 || str[i-1] != '\\' {
				idx = i
				break
			}
		}
	}
	return idx
}

// IsAlphanum returns true if the string contains only
// alphanumeric characters or the underscore character,
// false otherwise.
// @param c the string in question
// @return true if c is alphanumeric or underscore, false if not
func IsAlphanum(s string) bool {
	for _, r := range s {
		if !IsAlpha(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

// IsAlpha returns true if the rune contains only 'a'..'z' and 'A'..'Z'.
func IsAlpha(r rune) bool {
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
		return false
	}
	return true
}

// ValidateURI is not part of the reference implementation pseudo code
// found in the CPE 2.3 specification.  It enforces two rules in the
// specification:
//   URI must start with the characters "cpe:/"
//   A URI may not contain more than 7 components
// If either rule is violated, a ParseErr is thrown.
func ValidateURI(in string) error {
	in = strings.ToLower(in)
	// make sure uri starts with cpe:/
	if !strings.HasPrefix(in, "cpe:/") {
		return errors.Wrapf(ErrParse, "Error: URI must start with 'cpe:/'.  Given: %s", in)
	}
	// make sure uri doesn't contain more than 7 colons
	if count := strings.Count(in, ":"); count > 7 {
		return errors.Wrapf(ErrParse, "Error parsing URI.  Found %d extra components in: %s", count-7, in)
	}
	return nil
}

// ValidateFS is not part of the reference implementation pseudo code
// found in the CPE 2.3 specification.  It enforces three rules found in the
// specification:
//    Formatted string must start with the characters "cpe:2.3:"
//    A formatted string must contain 11 components
//    A formatted string must not contain empty components
// If any rule is violated, a ParseException is thrown.
func ValidateFS(in string) error {
	in = strings.ToLower(in)
	if !strings.HasPrefix(in, "cpe:2.3:") {
		return errors.Wrapf(ErrParse, "Error: Formatted String must start with \"cpe:2.3\". Given: %s", in)
	}
	// make sure fs contains exactly 12 unquoted colons
	count := 0
	for i := 0; i != len(in); i++ {
		if in[i] == ':' {
			if i == 0 || in[i-1] != '\\' {
				count++
			} else {
				continue
			}
			if i < len(in)-1 && in[i+1] == ':' {
				return errors.Wrap(ErrParse, "Error parsing formatted string. Found empty component")
			}
		}
	}
	if count > 12 {
		extra := count - 12
		return errors.Wrapf(ErrParse, "Error parsing formatted string. Found %d components in: %s", extra, in)
	}
	if count < 12 {
		missing := 12 - count
		return errors.Wrapf(ErrParse, "Error parsing formatted string. Missing %d components in: %s", missing, in)
	}
	return nil
}
