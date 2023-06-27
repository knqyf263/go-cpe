package naming

import (
	"strings"

	"github.com/knqyf263/go-cpe/common"
	"github.com/pkg/errors"
)

// UnbindURI is a top level function used to unbind a URI to a WFN.
// @param uri String representing the URI to be unbound.
// @return WellFormedName representing the unbound URI.
func UnbindURI(uri string) (common.WellFormedName, error) {
	if err := common.ValidateURI(uri); err != nil {
		return nil, errors.Wrap(err, "Failed to validate uri")
	}
	result := common.WellFormedName{}
	for i := 0; i < 8; i++ {
		v := getCompURI(uri, i)
		d, err := decode(v)
		if i != 6 && err != nil {
			return nil, err
		}
		switch i {
		case 1:
			err = result.Set(common.AttributePart, d)
		case 2:
			err = result.Set(common.AttributeVendor, d)
		case 3:
			err = result.Set(common.AttributeProduct, d)
		case 4:
			err = result.Set(common.AttributeVersion, d)
		case 5:
			err = result.Set(common.AttributeUpdate, d)
		case 6:
			// Special handling for edition component.
			// Unpack edition if needed.
			if v == "" || v == "-" || v[0] != '~' {
				// Just a logical value or a non-packed value.
				// So unbind to legacy edition, leaving other
				// extended attributes unspecified.
				err = result.Set(common.AttributeEdition, d)
			} else {
				// We have five values packed together here.
				if result, err = unpack(v, result); err != nil {
					return nil, err
				}
			}
		case 7:
			err = result.Set(common.AttributeLanguage, d)
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil

}

// UnbindFS is a top level function to unbind a formatted string to WFN.
// @param fs Formatted string to unbind
// @return WellFormedName
// @throws ParseException
func UnbindFS(fs string) (common.WellFormedName, error) {
	// Validate the formatted string
	if err := common.ValidateFS(fs); err != nil {
		return nil, err
	}
	// Initialize empty WFN
	result := common.NewWellFormedName()
	// The cpe scheme is the 0th component, the cpe version is the 1st.
	// So we start parsing at the 2nd component.
	for a := 2; a != 13; a++ {
		// Get the a'th string field.
		s := getCompFS(fs, a)
		// Unbind the string.
		v, err := unbindValueFS(s)
		if err != nil {
			return nil, err
		}

		// Set the value of the corresponding attribute.
		switch a {
		case 2:
			err = result.Set(common.AttributePart, v)
			break
		case 3:
			err = result.Set(common.AttributeVendor, v)
			break
		case 4:
			err = result.Set(common.AttributeProduct, v)
			break
		case 5:
			err = result.Set(common.AttributeVersion, v)
			break
		case 6:
			err = result.Set(common.AttributeUpdate, v)
			break
		case 7:
			err = result.Set(common.AttributeEdition, v)
			break
		case 8:
			err = result.Set(common.AttributeLanguage, v)
			break
		case 9:
			err = result.Set(common.AttributeSwEdition, v)
			break
		case 10:
			err = result.Set(common.AttributeTargetSw, v)
			break
		case 11:
			err = result.Set(common.AttributeTargetHw, v)
			break
		case 12:
			err = result.Set(common.AttributeOther, v)
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// getCompURI return the i'th component of the URI.
// @param uri String representation of URI to retrieve components from.
// @param i Index of component to return.
// @return If i = 0, returns the URI scheme. Otherwise, returns the i'th
// component of uri.
func getCompURI(uri string, i int) string {
	if i == 0 {
		return uri[:strings.Index(uri, "/")]
	}
	sa := strings.Split(uri, ":")
	// If requested component exceeds the number
	// of components in URI, return blank
	if i >= len(sa) {
		return ""
	}
	if i == 1 {
		return strings.TrimLeft(sa[1], "/")
	}
	return sa[i]
}

// Returns the i'th field of the formatted string.  The colon is the field
// delimiter unless prefixed by a backslash.
// @param fs formatted string to retrieve from
// @param i index of field to retrieve from fs.
// @return value of index of formatted string
func getCompFS(fs string, i int) string {
	if i < 0 {
		return ""
	}

	for j := 0; j < i; j++ {
		// return the substring from index 0 to the first occurence of an
		// unescaped colon
		colonIdx := common.GetUnescapedColonIndex(fs)
		if colonIdx == 0 {
			fs = ""
			break
		}
		fs = fs[colonIdx+1:]
	}
	endIdx := common.GetUnescapedColonIndex(fs)
	// If no colon is found, we are at the end of the formatted string,
	// so just return what's left.
	if endIdx == 0 {
		return fs
	}
	return fs[:endIdx]
}

// Takes a string value and returns the appropriate logical value if string
// is the bound form of a logical value.  If string is some general value
// string, add quoting of non-alphanumerics as needed.
// @param s value to be unbound
// @return logical value or quoted string
// @throws ParseException
func unbindValueFS(s string) (interface{}, error) {
	if s == "*" {
		any, _ := common.NewLogicalValue("ANY")
		return any, nil
	}
	if s == "-" {
		na, _ := common.NewLogicalValue("NA")
		return na, nil
	}
	result, err := addQuoting(s)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Inspect each character in a string, copying quoted characters, with
// their escaping, into the result.  Look for unquoted non alphanumerics
// and if not "*" or "?", add escaping.
// @param s
// @return
// @throws ParseException
func addQuoting(s string) (result string, err error) {
	idx := 0
	embedded := false
	for idx < len(s) {
		c := s[idx : idx+1]
		if common.IsAlphanum(c) || c == "_" {
			// Alphanumeric characters pass untouched.
			result += c
			idx++
			embedded = true
			continue
		}
		if c == "\\" {
			// Anything quoted in the bound string stays quoted in the
			// unbound string.
			result += s[idx : idx+2]
			idx += 2
			embedded = true
			continue
		}
		if c == "*" {
			// An unquoted asterisk must appear at the beginning or the end
			// of the string.
			if idx == 0 || idx == len(s)-1 {
				result += c
				idx = idx + 1
				embedded = true
				continue
			}
			return "", errors.Wrap(common.ErrParse, "Error! cannot have unquoted * embedded in formatted string.")
		}
		if c == "?" {
			// An unquoted question mark must appear at the beginning or
			// end of the string, or in a leading or trailing sequence.
			// if  // ? legal at beginning or end
			valid := false
			if idx == 0 || idx == len(s)-1 {
				// ? legal at beginning or end
				valid = true
			} else if !embedded && idx > 0 && s[idx-1:idx] == "?" {
				// embedded is false, so must be preceded by ?
				valid = true
			} else if embedded && len(s) >= idx+2 && s[idx+1:idx+2] == "?" {
				// embedded is true, so must be followed by ?
				valid = true
			}

			if !valid {
				return "", errors.Wrap(common.ErrParse, "Error! cannot have unquoted ? embedded in formatted string.")
			}
			result += c
			idx++
			embedded = false
			continue
		}
		// All other characters must be quoted.
		result += "\\" + c
		idx++
		embedded = true
	}
	return result, nil
}

// decode scans a string and returns a copy with all percent-encoded characters
// decoded.  This function is the inverse of pctEncode() defined in the
// CPENameBinder class.  Only legal percent-encoded forms are decoded.
// Others raise a ParseException.
// @param s String to be decoded
// @return decoded string
// @throws ParseException
// @see CPENameBinder#pctEncode(java.lang.String)
func decode(s string) (interface{}, error) {
	if s == "" {
		any, _ := common.NewLogicalValue("ANY")
		return any, nil
	}
	if s == "-" {
		na, _ := common.NewLogicalValue("NA")
		return na, nil
	}
	// Start the scanning loop.
	// Normalize: convert all uppercase letters to lowercase first.
	s = strings.ToLower(s)
	result := ""
	idx := 0
	embedded := false
	for idx < len(s) {
		// Get the idx'th character of s.
		c := s[idx : idx+1]
		// Deal with dot, hyphen, and tilde: decode with quoting.
		if c == "." || c == "-" || c == "~" {
			result += "\\" + c
			idx++
			// a non-%01 encountered.
			embedded = true
			continue
		}
		if c != "%" {
			result += c
			idx++
			// a non-%01 encountered.
			embedded = true
			continue
		}
		// We get here if we have a substring starting w/ '%'.
		form := s[idx : idx+3]
		if form == "%01" {
			valid := false
			if idx == 0 || idx == len(s)-3 {
				valid = true
			} else if !embedded && (idx-3 >= 0) && s[idx-3:idx] == "%01" {
				valid = true
			} else if embedded && len(s) >= idx+6 && s[idx+3:idx+6] == "%01" {
				valid = true
			}

			if valid {
				result += "?"
				idx = idx + 3
				continue
			} else {
				return nil, errors.Wrapf(common.ErrParse, "Error decoding string")
			}
		} else if form == "%02" {
			if idx == 0 || idx == len(s)-3 {
				result += "*"
			} else {
				return nil, errors.Wrapf(common.ErrParse, "Error decoding string")
			}
		} else if form == "%21" {
			result += `\!`
		} else if form == "%22" {
			result += `\"`
		} else if form == "%23" {
			result += `\#`
		} else if form == "%24" {
			result += `\$`
		} else if form == "%25" {
			result += `\%`
		} else if form == "%26" {
			result += `\&`
		} else if form == "%27" {
			result += `\'`
		} else if form == "%28" {
			result += `\(`
		} else if form == "%29" {
			result += `\)`
		} else if form == "%2a" {
			result += `\*`
		} else if form == "%2b" {
			result += `\+`
		} else if form == "%2c" {
			result += `\,`
		} else if form == "%2f" {
			result += `\/`
		} else if form == "%3a" {
			result += `\:`
		} else if form == "%3b" {
			result += `\;`
		} else if form == "%3c" {
			result += `\<`
		} else if form == "%3d" {
			result += `\=`
		} else if form == "%3e" {
			result += `\>`
		} else if form == "%3f" {
			result += `\?`
		} else if form == "%40" {
			result += `\@`
		} else if form == "%5b" {
			result += `\[`
		} else if form == "%5c" {
			result += `\\`
		} else if form == "%5d" {
			result += `\]`
		} else if form == "%5e" {
			result += `\^`
		} else if form == "%60" {
			result += "\\`"
		} else if form == "%7b" {
			result += `\{`
		} else if form == "%7c" {
			result += `\|`
		} else if form == "%7d" {
			result += `\}`
		} else if form == "%7e" {
			result += `\~`
		} else {
			return nil, errors.Wrapf(common.ErrParse, "Unknown form: %s", form)
		}
		idx = idx + 3
		embedded = true
	}
	return result, nil
}

// unpack unpacks the elements in s and sets the attributes in the given
// WellFormedName accordingly.
// @param s packed String
// @param wfn WellFormedName
// @return The augmented WellFormedName.
func unpack(s string, wfn common.WellFormedName) (common.WellFormedName, error) {
	// Set each component in the WFN.
	editions := strings.SplitN(s[1:], "~", 5)
	if len(editions) != 5 {
		return nil, errors.Wrap(common.ErrParse, "editions must be 5")
	}
	attributes := []string{common.AttributeEdition, common.AttributeSwEdition, common.AttributeTargetSw,
		common.AttributeTargetHw, common.AttributeOther}
	for i, a := range attributes {
		e, err := decode(editions[i])
		if err != nil {
			return nil, err
		}
		if err = wfn.Set(a, e); err != nil {
			return nil, err
		}
	}
	return wfn, nil
}
