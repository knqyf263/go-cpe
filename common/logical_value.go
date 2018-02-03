package common

import (
	"errors"
	"strings"
)

var (
	// ErrIllegalArgument is returned when illegal argument
	ErrIllegalArgument = errors.New("Illegal argument")
)

// LogicalValue represents a Logical Value.
// @see <a href="http://cpe.mitre.org">cpe.mitre.org</a> for more information.
// @author JKRAUNELIS
// @email jkraunelis@mitre.org
type LogicalValue struct {
	Any bool
	Na  bool
}

// NewLogicalValue returns Logicalvalue
func NewLogicalValue(t string) (lv LogicalValue, err error) {
	t = strings.ToUpper(t)
	if t == "ANY" {
		lv.Any = true
	} else if t == "NA" {
		lv.Na = true
	} else {
		return LogicalValue{}, ErrIllegalArgument
	}
	return lv, nil
}

// IsANY returns whether any is true
func (lv LogicalValue) IsANY() bool {
	return lv.Any
}

// IsNA returns whether na is true
func (lv LogicalValue) IsNA() bool {
	return lv.Na
}

// String : String
func (lv LogicalValue) String() string {
	if lv.Any {
		return "ANY"
	}
	return "NA"
}
