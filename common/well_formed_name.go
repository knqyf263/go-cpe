package common

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

const (
	AttributePart      = "part"
	AttributeVendor    = "vendor"
	AttributeProduct   = "product"
	AttributeVersion   = "version"
	AttributeUpdate    = "update"
	AttributeEdition   = "edition"
	AttributeLanguage  = "language"
	AttributeSwEdition = "sw_edition"
	AttributeTargetSw  = "target_sw"
	AttributeTargetHw  = "target_hw"
	AttributeOther     = "other"
)

var (
	attributes = []string{AttributePart, AttributeVendor, AttributeProduct, AttributeVersion, AttributeUpdate,
		AttributeEdition, AttributeLanguage, AttributeSwEdition, AttributeTargetSw, AttributeTargetHw, AttributeOther}

	// ErrIllegalAttribute is returned when illegal argument
	ErrIllegalAttribute = errors.New("Illegal attribute")
	// ErrParse is returned when parse error
	ErrParse = errors.New("Parse error")
)

// WellFormedName represents a Well Formed Name, as defined
// in the CPE Specification version 2.3.
//
// @see <a href="http://cpe.mitre.org">cpe.mitre.org</a> for details.
type WellFormedName map[string]interface{}

// NewWellFormedName constructs a new WellFormedName object, with all components set to the default value "ANY".
func NewWellFormedName() WellFormedName {
	wfn := WellFormedName{}
	for _, a := range attributes {
		if a != AttributePart {
			wfn[a], _ = NewLogicalValue("ANY")
		}
	}
	return wfn
}

// Initialize sets each component to the given parameter value.
// If a parameter is null, the component is set to the default value "ANY".
// @param part string representing the part component
// @param vendor string representing the vendor component
// @param product string representing the product component
// @param version string representing the version component
// @param update string representing the update component
// @param edition string representing the edition component
// @param language string representing the language component
// @param sw_edition string representing the sw_edition component
// @param target_sw string representing the target_sw component
// @param target_hw string representing the target_hw component
// @param other string representing the other component
func (wfn WellFormedName) Initialize(part, vendor, product, version, update, edition, language, swEdition, targetSw, targetHw, other interface{}) {
	wfn[AttributePart] = part
	wfn[AttributeVendor] = vendor
	wfn[AttributeProduct] = product
	wfn[AttributeVersion] = version
	wfn[AttributeUpdate] = update
	wfn[AttributeEdition] = edition
	wfn[AttributeLanguage] = language
	wfn[AttributeSwEdition] = swEdition
	wfn[AttributeTargetSw] = targetSw
	wfn[AttributeTargetHw] = targetHw
	wfn[AttributeOther] = other
}

// Get gets attribute
// @param attribute String representing the component value to get
// @return the String value of the given component, or default value "ANY"
// if the component does not exist
func (wfn WellFormedName) Get(attribute string) interface{} {
	if v, ok := wfn[attribute]; ok {
		return v
	}
	any, _ := NewLogicalValue("ANY")
	return any
}

// Set sets the given attribute to value, if the attribute is in the list of permissible components
// @param attribute String representing the component to set
// @param value Object representing the value of the given component
func (wfn WellFormedName) Set(attribute string, value interface{}) (err error) {
	if valid := IsValidAttribute(attribute); !valid {
		return ErrIllegalAttribute
	}

	if value == nil {
		wfn[attribute], _ = NewLogicalValue("ANY")
		return nil
	}

	if _, ok := value.(LogicalValue); ok {
		if attribute == AttributePart {
			return errors.Wrap(ErrIllegalAttribute, "part component cannot be a logical value")
		}
		wfn[attribute] = value
		return nil
	}

	svalue, ok := value.(string)
	if !ok {
		return errors.Wrap(ErrIllegalAttribute, "value must be a logical value or string")
	}

	if err = ValidateStringValue(svalue); err != nil {
		return errors.Wrap(err, "Failed to validate a value")
	}

	// part must be a, o, or h
	if attribute == AttributePart {
		if svalue != "a" && svalue != "o" && svalue != "h" {
			return errors.Wrapf(ErrParse, "part component must be one of the following: 'a', 'o', 'h': %s", svalue)
		}
	}

	// should be good to go
	wfn[attribute] = value

	return nil
}

// GetString gets attribute as string
// @param attribute String representing the component value to get
// @return the String value of the given component, or default value "ANY"
// if the component does not exist
func (wfn WellFormedName) GetString(attribute string) string {
	return fmt.Sprintf("%s", wfn.Get(attribute))
}

// String returns string representation of the WellFormedName
func (wfn WellFormedName) String() string {
	s := "wfn:["
	for _, attr := range attributes {
		s += fmt.Sprintf("%s=", attr)
		o := wfn.Get(attr)
		if lv, ok := o.(LogicalValue); ok {
			s += fmt.Sprintf("%s, ", lv)
		} else {
			s += fmt.Sprintf("\"%s\", ", o)
		}
	}
	s = strings.TrimSpace(s)
	s = strings.TrimRight(s, ",")
	s += "]"
	return s
}

// IsValidAttribute validates an attribute name
func IsValidAttribute(attribute string) (valid bool) {
	for _, a := range attributes {
		if a == attribute {
			valid = true
		}
	}
	return valid
}

// ValidateStringValue validates an string value
func ValidateStringValue(svalue string) (err error) {
	// svalue has more than one unquoted star
	if strings.HasPrefix(svalue, "**") || strings.HasSuffix(svalue, "**") {
		return errors.Wrapf(ErrParse, "component cannot contain more than one * in sequence: %s", svalue)
	}

	prev := ' ' // dummy value
	for i, r := range svalue {
		// check for printable characters - no control characters
		if !unicode.IsPrint(r) {
			return errors.Wrapf(ErrParse, "encountered non printable character in: %s", svalue)
		}
		// svalue has whitespace
		if unicode.IsSpace(r) {
			return errors.Wrapf(ErrParse, "component cannot contain whitespace:: %s", svalue)
		}
		if unicode.IsPunct(r) && prev != '\\' && r != '\\' {
			// svalue has an unquoted *
			if r == '*' && (i != 0 && i != len(svalue)-1) {
				return errors.Wrapf(ErrParse, "component cannot contain embedded *: %s", svalue)
			}

			if r != '*' && r != '?' && r != '_' {
				// svalue has unquoted punctuation embedded
				return errors.Wrapf(ErrParse, "component cannot contain unquoted punctuation: %s", svalue)
			}
		}
		prev = r
	}

	if strings.Contains(svalue, "?") {
		if svalue == "?" {
			// single ? is valid
			return nil
		}

		s := strings.Trim(svalue, "?")
		if ContainsQuestions(s) {
			return errors.Wrapf(ErrParse, "component cannot contain embedded ?: %s", svalue)
		}
	}

	// single asterisk is not allowed
	if svalue == "*" {
		return errors.Wrapf(ErrParse, "component cannot be a single *: %s", svalue)
	}

	// quoted hyphen not allowed by itself
	if svalue == `\-` {
		return errors.Wrapf(ErrParse, "component cannot be quoted hyphen: %s", svalue)
	}

	return nil
}
