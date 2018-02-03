package common

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

var (
	any, _ = NewLogicalValue("ANY")
	na, _  = NewLogicalValue("NA")
)

func TestNewWellFormedName(t *testing.T) {
	vectors := []struct {
		expected WellFormedName
	}{{
		expected: WellFormedName{
			"vendor":     any,
			"product":    any,
			"version":    any,
			"update":     any,
			"edition":    any,
			"language":   any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
		},
	},
	}

	for i, v := range vectors {
		actual := NewWellFormedName()
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("test %d, Result: %v, want %v", i, actual, v.expected)
		}
	}
}

func TestWellFormedNameInitialize(t *testing.T) {
	vectors := []struct {
		part      interface{}
		vendor    interface{}
		product   interface{}
		version   interface{}
		update    interface{}
		edition   interface{}
		language  interface{}
		swEdition interface{}
		targetSw  interface{}
		targetHw  interface{}
		other     interface{}
		expected  WellFormedName
	}{{
		part:      "a",
		vendor:    "microsoft",
		product:   "windows_7",
		version:   na,
		update:    any,
		edition:   any,
		language:  any,
		swEdition: any,
		targetSw:  any,
		targetHw:  any,
		other:     any,
		expected: WellFormedName{
			"part":       "a",
			"vendor":     "microsoft",
			"product":    "windows_7",
			"version":    na,
			"update":     any,
			"edition":    any,
			"language":   any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
		},
	},
	}

	for i, v := range vectors {
		wfn := WellFormedName{}
		wfn.Initialize(v.part, v.vendor, v.product, v.version, v.update, v.edition, v.language,
			v.swEdition, v.targetSw, v.targetHw, v.other)
		if !reflect.DeepEqual(wfn, v.expected) {
			t.Errorf("test %d, Result: %v, want %v", i, wfn, v.expected)
		}
	}
}

func TestWellFormedNameGet(t *testing.T) {
	vectors := []struct {
		wfn       WellFormedName
		attribute string
		expected  interface{}
	}{{
		wfn: WellFormedName{
			"part": "a",
		},
		attribute: "part",
		expected:  "a",
	}, {
		wfn: WellFormedName{
			"vendor": "microsoft",
		},
		attribute: "vendor",
		expected:  "microsoft",
	}, {
		wfn: WellFormedName{
			"product": any,
		},
		attribute: "product",
		expected:  any,
	}, {
		wfn: WellFormedName{
			"product": any,
		},
		attribute: "version",
		expected:  any,
	},
	}

	for i, v := range vectors {
		actual := v.wfn.Get(v.attribute)
		if !reflect.DeepEqual(actual, v.expected) {
			t.Errorf("test %d, Result: %v, want %v", i, actual, v.expected)
		}
	}
}

func TestWellFormedNameSet(t *testing.T) {
	vectors := []struct {
		wfn       WellFormedName
		attribute string
		value     interface{}
		expected  WellFormedName
		wantErr   error
	}{{
		wfn:       WellFormedName{},
		attribute: "part",
		value:     "a",
		expected: WellFormedName{
			"part": "a",
		},
	}, {
		wfn: WellFormedName{
			"part": "a",
		},
		attribute: "vendor",
		value:     any,
		expected: WellFormedName{
			"part":   "a",
			"vendor": any,
		},
	}, {
		wfn:       WellFormedName{},
		attribute: "vendor",
		value:     nil,
		expected: WellFormedName{
			"vendor": any,
		},
	}, {
		wfn:       WellFormedName{},
		attribute: "part",
		value:     na, // part component cannot be a logical value
		wantErr:   ErrIllegalAttribute,
	}, {
		wfn:       WellFormedName{},
		attribute: "part",
		value:     "i", // part component must be one of the following: 'a', 'o', 'h'
		wantErr:   ErrParse,
	}, {
		wfn:       WellFormedName{},
		attribute: "vendor",
		value:     1, // value must be a logical value or string
		wantErr:   ErrIllegalAttribute,
	}, {
		wfn:       WellFormedName{},
		attribute: "invalid", // invalid attribute name
		value:     "invalid",
		wantErr:   ErrIllegalAttribute,
	}, {
		wfn:       WellFormedName{},
		attribute: "version",
		value:     "**1.2.3",
		wantErr:   ErrParse,
	},
	}

	for i, v := range vectors {
		err := v.wfn.Set(v.attribute, v.value)
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %d, Error: %v, want %v", i, errors.Cause(err), v.wantErr)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(v.wfn, v.expected) {
			t.Errorf("test %d, Result: %v, want %v", i, v.wfn, v.expected)
		}
	}
}

func TestWellFormedNameString(t *testing.T) {
	vectors := []struct {
		wfn      WellFormedName
		expected string
	}{{
		wfn:      WellFormedName{},
		expected: `wfn:[part=ANY, vendor=ANY, product=ANY, version=ANY, update=ANY, edition=ANY, language=ANY, sw_edition=ANY, target_sw=ANY, target_hw=ANY, other=ANY]`,
	}, {
		wfn: WellFormedName{
			"part":       "a",
			"vendor":     "microsoft",
			"product":    "windows_7",
			"version":    na,
			"update":     any,
			"edition":    any,
			"language":   any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
		},
		expected: `wfn:[part="a", vendor="microsoft", product="windows_7", version=NA, update=ANY, edition=ANY, language=ANY, sw_edition=ANY, target_sw=ANY, target_hw=ANY, other=ANY]`,
	},
	}

	for i, v := range vectors {
		actual := v.wfn.String()
		if actual != v.expected {
			t.Errorf("test %d, Result: %v, want %v", i, actual, v.expected)
		}
	}
}

func TestValidateStringValue(t *testing.T) {
	vectors := []struct {
		svalue  string
		wantErr error
	}{{
		svalue: "foo",
	}, {
		svalue: "?",
	}, {
		svalue: "??",
	}, {
		svalue:  "**foo",
		wantErr: ErrParse,
	}, {
		svalue:  "foo**",
		wantErr: ErrParse,
	}, {
		svalue:  "foo*bar",
		wantErr: ErrParse,
	}, {
		svalue: "*foo*",
	}, {
		svalue:  "foo!bar",
		wantErr: ErrParse,
	}, {
		svalue:  "foo/bar",
		wantErr: ErrParse,
	}, {
		svalue:  "foo\x07bar",
		wantErr: ErrParse,
	}, {
		svalue:  "foo bar",
		wantErr: ErrParse,
	}, {
		svalue: "???foo??",
	}, {
		svalue:  "??foo?bar??",
		wantErr: ErrParse,
	}, {
		svalue:  "*",
		wantErr: ErrParse,
	}, {
		svalue:  `\-`,
		wantErr: ErrParse,
	},
	}

	for i, v := range vectors {
		err := ValidateStringValue(v.svalue)
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %d, Error: %v, want %v", i, errors.Cause(err), v.wantErr)
		}
	}
}
