package naming

import (
	"fmt"
	"strings"

	"github.com/knqyf263/go-cpe/common"
)

// BindToURI a {@link WellFormedName} object to a URI.
// @param w WellFormedName to be bound to URI
// @return URI binding of WFN
func BindToURI(w common.WellFormedName) (uri string) {
	// Initialize the output with the CPE v2.2 URI prefix.
	uri = "cpe:/"

	// Define the attributes that correspond to the seven components in a v2.2. CPE.
	attributes := []string{"part", "vendor", "product", "version", "update", "edition", "language"}

	for _, a := range attributes {
		v := ""
		if a == "edition" {
			// Call the pack() helper function to compute the proper binding for the edition element.
			ed := bindValueForURI(w.Get(common.AttributeEdition))
			swEd := bindValueForURI(w.Get(common.AttributeSwEdition))
			targetSw := bindValueForURI(w.Get(common.AttributeTargetSw))
			targetHw := bindValueForURI(w.Get(common.AttributeTargetHw))
			other := bindValueForURI(w.Get(common.AttributeOther))
			v = pack(ed, swEd, targetSw, targetHw, other)
		} else {
			// Get the value for a in w, then bind to a string
			// for inclusion in the URI.
			v = bindValueForURI(w.Get(a))
		}
		// Append v to the URI then add a colon.
		uri += v + ":"
	}
	return strings.TrimRight(uri, ":")
}

// BindToFS is top-level function used to bind WFN w to formatted string.
// @param w WellFormedName to bind
// @return Formatted String
func BindToFS(w common.WellFormedName) (fs string) {
	// Initialize the output with the CPE v2.3 string prefix.
	fs = "cpe:2.3:"
	attributes := []string{"part", "vendor", "product", "version",
		"update", "edition", "language", "sw_edition", "target_sw",
		"target_hw", "other"}
	for _, a := range attributes {
		v := bindValueForFS(w.Get(a))
		fs += v
		if a != common.AttributeOther {
			fs += ":"
		}
	}
	return fs
}

// bindValueForURI converts a string to the proper string for including in a CPE v2.2-conformant URI.
// The logical value ANY binds to the blank in the 2.2-conformant URI.
// @param s string to be converted
// @return converted string
func bindValueForURI(s interface{}) string {
	if lv, ok := s.(common.LogicalValue); ok {
		// The value NA binds to a blank.
		if lv.IsANY() {
			return ""
		}
		// The value NA binds to a single hyphen.
		if lv.IsNA() {
			return "-"
		}
	}

	if str, ok := s.(string); ok {
		return transformForURI(str)
	}
	return ""
}

// bindValueForFS converts the value v to its proper string representation for insertion to formatted string.
// @param v value to convert
// @return Formatted value
func bindValueForFS(v interface{}) string {
	if lv, ok := v.(common.LogicalValue); ok {
		// The value NA binds to a asterisk.
		if lv.IsANY() {
			return "*"
		}
		// The value NA binds to a single hyphen.
		if lv.IsNA() {
			return "-"
		}
	}
	if str, ok := v.(string); ok {
		return processQuotedChars(str)
	}
	return ""
}

// Inspect each character in string s.  Certain nonalpha characters pass
// thru without escaping into the result, but most retain escaping.
// @param s
// @return
func processQuotedChars(s string) (result string) {
	idx := 0
	for idx < len(s) {
		c := s[idx : idx+1]
		if c == "\\" {
			// escaped characters are examined.
			nextchr := s[idx+1 : idx+2]
			// the period, hyphen and underscore pass unharmed.
			if nextchr == "." || nextchr == "-" || nextchr == "_" {
				result += nextchr
				idx += 2
				continue
			} else {
				// all others retain escaping.
				result += "\\" + nextchr
				idx += 2
				continue
			}
		}
		// unquoted characters pass thru unharmed.
		result += c
		idx = idx + 1
	}
	return result
}

// transformForURI scans an input string and performs the following transformations:
// - Pass alphanumeric characters thru untouched
// - Percent-encode quoted non-alphanumerics as needed
// - Unquoted special characters are mapped to their special forms
// @param s string to be transformed
// @return transformed string
func transformForURI(s string) (result string) {
	idx := 0
	for idx < len(s) {
		// Get the idx'th character of s.
		thischar := s[idx : idx+1]
		if common.IsAlphanum(thischar) {
			result += thischar
			idx++
			continue
		}
		// Check for escape character.
		if thischar == "\\" {
			idx++
			nxtchar := s[idx : idx+1]
			result += pctEncode(nxtchar)
			idx++
			continue
		}
		// Bind the unquoted '?' special character to "%01".
		if thischar == "?" {
			result += "%01"
		}
		// Bind the unquoted '*' special character to "%02".
		if thischar == "*" {
			result += "%02"
		}
		idx++

	}
	return result
}

// pctEncode returns the appropriate percent-encoding of character c.
// Certain characters are returned without encoding.
// @param c the single character string to be encoded
// @return the percent encoded string
func pctEncode(c string) string {
	if c == "!" {
		return "%21"
	}
	if c == "\"" {
		return "%22"
	}
	if c == "#" {
		return "%23"
	}
	if c == "$" {
		return "%24"
	}
	if c == "%" {
		return "%25"
	}
	if c == "&" {
		return "%26"
	}
	if c == "'" {
		return "%27"
	}
	if c == "(" {
		return "%28"
	}
	if c == ")" {
		return "%29"
	}
	if c == "*" {
		return "%2a"
	}
	if c == "+" {
		return "%2b"
	}
	if c == "," {
		return "%2c"
	}
	// bound without encoding.
	if c == "-" {
		return c
	}
	// bound without encoding.
	if c == "." {
		return c
	}
	if c == "/" {
		return "%2f"
	}
	if c == ":" {
		return "%3a"
	}
	if c == ";" {
		return "%3b"
	}
	if c == "<" {
		return "%3c"
	}
	if c == "=" {
		return "%3d"
	}
	if c == ">" {
		return "%3e"
	}
	if c == "?" {
		return "%3f"
	}
	if c == "@" {
		return "%40"
	}
	if c == "[" {
		return "%5b"
	}
	if c == "\\" {
		return "%5c"
	}
	if c == "]" {
		return "%5d"
	}
	if c == "^" {
		return "%5e"
	}
	if c == "`" {
		return "%60"
	}
	if c == "{" {
		return "%7b"
	}
	if c == "|" {
		return "%7c"
	}
	if c == "}" {
		return "%7d"
	}
	if c == "~" {
		return "%7e"
	}
	// Shouldn't reach here, return original character
	return c
}

/**
 * Packs the values of the five arguments into the single
 * edition component.  If all the values are blank, the
 * function returns a blank.
 * @param ed edition string
 * @param swEd software edition string
 * @param tSw target software string
 * @param tHw target hardware string
 * @param oth other edition information string
 * @return the packed string, or blank
 */
func pack(ed, swEd, tSw, tHw, oth string) string {
	if swEd == "" && tSw == "" && tHw == "" && oth == "" {
		// All the extended attributes are blank, so don't do
		// any packing, just return ed.
		return ed
	}
	// Otherwise, pack the five values into a single string
	// prefixed and internally delimited with the tilde.
	return fmt.Sprintf("~%s~%s~%s~%s~%s", ed, swEd, tSw, tHw, oth)
}
