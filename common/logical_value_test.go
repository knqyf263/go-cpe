package common

import (
	"testing"

	"github.com/pkg/errors"
)

func TestNewLogicalValue(t *testing.T) {
	vectors := []struct {
		s        string
		expected LogicalValue
		wantErr  error
	}{{
		s:        `ANY`,
		expected: LogicalValue{Any: true},
	}, {
		s:        `NA`,
		expected: LogicalValue{Na: true},
	}, {
		s:        `any`,
		expected: LogicalValue{Any: true},
	}, {
		s:        `na`,
		expected: LogicalValue{Na: true},
	}, {
		s:       `invalid`,
		wantErr: ErrIllegalArgument,
	},
	}

	for i, v := range vectors {
		actual, err := NewLogicalValue(v.s)
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %d, Error: got %v, want %v", i, errors.Cause(err), v.wantErr)
		}
		if err != nil {
			continue
		}
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestIsANY(t *testing.T) {
	vectors := []struct {
		lv       LogicalValue
		expected bool
	}{{
		lv:       LogicalValue{Any: true, Na: false},
		expected: true,
	}, {
		lv:       LogicalValue{Any: false, Na: true},
		expected: false,
	},
	}

	for i, v := range vectors {
		actual := v.lv.IsANY()
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestIsNA(t *testing.T) {
	vectors := []struct {
		lv       LogicalValue
		expected bool
	}{{
		lv:       LogicalValue{Any: true, Na: false},
		expected: false,
	}, {
		lv:       LogicalValue{Any: false, Na: true},
		expected: true,
	},
	}

	for i, v := range vectors {
		actual := v.lv.IsNA()
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestLogicalValueString(t *testing.T) {
	vectors := []struct {
		lv       LogicalValue
		expected string
	}{{
		lv:       LogicalValue{Any: true, Na: false},
		expected: "ANY",
	}, {
		lv:       LogicalValue{Any: false, Na: true},
		expected: "NA",
	},
	}

	for i, v := range vectors {
		actual := v.lv.String()
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}
