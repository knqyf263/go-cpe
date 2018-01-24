package common

import (
	"testing"

	"github.com/pkg/errors"
)

func TestCountEscapeCharacters(t *testing.T) {
	vectors := []struct {
		s        string
		expected int
	}{{
		s:        `\abc`,
		expected: 1,
	}, {
		s:        `\\abc`,
		expected: 1,
	}, {
		s:        `\\\abc`,
		expected: 2,
	}, {
		s:        `\\\\abc`,
		expected: 2,
	}, {
		s:        `\\abc\d`,
		expected: 2,
	}, {
		s:        `\\abc\\d`,
		expected: 2,
	}, {
		s:        `\\abc\\\d`,
		expected: 3,
	},
	}

	for i, v := range vectors {
		actual := CountEscapeCharacters(v.s)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestLengthWithEscapeCharacters(t *testing.T) {
	vectors := []struct {
		s        string
		expected int
	}{{
		s:        `\abc`,
		expected: 3,
	}, {
		s:        `\\abc`,
		expected: 4,
	}, {
		s:        `\\\abc`,
		expected: 4,
	}, {
		s:        `\\\\abc`,
		expected: 5,
	}, {
		s:        `\\abc\d`,
		expected: 5,
	}, {
		s:        `\\abc\\d`,
		expected: 6,
	}, {
		s:        `\\abc\\\d`,
		expected: 6,
	},
	}

	for i, v := range vectors {
		actual := LengthWithEscapeCharacters(v.s)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestContainsWildCards(t *testing.T) {
	vectors := []struct {
		s        string
		expected bool
	}{{
		s:        "*abc",
		expected: true,
	}, {
		s:        "?abc",
		expected: true,
	}, {
		s:        "abc*",
		expected: true,
	}, {
		s:        `abc\*`,
		expected: false,
	}, {
		s:        "abc??",
		expected: true,
	}, {
		s:        "abc*def",
		expected: true,
	}, {
		s:        "abc?def",
		expected: true,
	}, {
		s:        `abc\:def*`,
		expected: true,
	}, {
		s:        `abc\:def:ghi`,
		expected: false,
	}, {
		s:        `abc\:def:ghi*`,
		expected: true,
	}, {
		s:        `abc\*def`,
		expected: false,
	}, {
		s:        `abc\*def*`,
		expected: true,
	},
	}

	for i, v := range vectors {
		actual := ContainsWildcards(v.s)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestContainsQuestions(t *testing.T) {
	vectors := []struct {
		s        string
		expected bool
	}{{
		s:        "*abc",
		expected: false,
	}, {
		s:        "?abc",
		expected: true,
	}, {
		s:        "abc*",
		expected: false,
	}, {
		s:        `abc\*`,
		expected: false,
	}, {
		s:        "abc??",
		expected: true,
	}, {
		s:        "abc*def",
		expected: false,
	}, {
		s:        "abc?def",
		expected: true,
	}, {
		s:        `abc\:def*`,
		expected: false,
	}, {
		s:        `abc\:def:ghi`,
		expected: false,
	}, {
		s:        `abc\:def:ghi*`,
		expected: false,
	}, {
		s:        `abc\?def`,
		expected: false,
	}, {
		s:        `abc\?def?`,
		expected: true,
	},
	}

	for i, v := range vectors {
		actual := ContainsQuestions(v.s)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestGetUnescapedColonIndex(t *testing.T) {
	vectors := []struct {
		s        string
		expected int
	}{{
		s:        "abc",
		expected: 0,
	}, {
		s:        ":abc",
		expected: 0,
	}, {
		s:        "abc:",
		expected: 3,
	}, {
		s:        `abc\:def`,
		expected: 0,
	}, {
		s:        `abc\:def:ghi`,
		expected: 8,
	},
	}

	for i, v := range vectors {
		actual := GetUnescapedColonIndex(v.s)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestIsAlphanum(t *testing.T) {
	vectors := []struct {
		s        string
		expected bool
	}{{
		s:        "abc",
		expected: true,
	}, {
		s:        "ABC123",
		expected: true,
	}, {
		s:        "xyz_789",
		expected: true,
	}, {
		s:        "abc XYZ",
		expected: false,
	}, {
		s:        "def!456",
		expected: false,
	}, {
		s:        "abcあABC",
		expected: false,
	},
	}

	for i, v := range vectors {
		actual := IsAlphanum(v.s)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	vectors := []struct {
		r        rune
		expected bool
	}{{
		r:        'a',
		expected: true,
	}, {
		r:        'A',
		expected: true,
	}, {
		r:        '!',
		expected: false,
	}, {
		r:        ' ',
		expected: false,
	}, {
		r:        'あ',
		expected: false,
	}, {},
	}

	for i, v := range vectors {
		actual := IsAlpha(v.r)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestValidateURI(t *testing.T) {
	vectors := []struct {
		s       string
		wantErr error
	}{{
		s: `cpe:/a:microsoft:internet_explorer:8.0.6001:beta::sp2`,
	}, {
		s:       `cpe:/a:microsoft:internet_explorer:8.0.6001:beta::sp2:invalid`,
		wantErr: ErrParse,
	}, {
		s:       `cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*:*:*`,
		wantErr: ErrParse,
	},
	}

	for i, v := range vectors {
		err := ValidateURI(v.s)
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %d, Error: got %v, want %v", i, err, v.wantErr)
		}
	}
}

func TestValidateFS(t *testing.T) {
	vectors := []struct {
		s       string
		wantErr error
	}{{
		s: `cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*:*:*`,
	}, {
		s: `cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*\::*:*`,
	}, {
		s: `cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*\:\::*:*`,
	}, {
		s:       "cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*:*:*:*",
		wantErr: ErrParse,
	}, {
		s:       "cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*::*",
		wantErr: ErrParse,
	}, {
		s:       "cpe:2.4:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*:*:*",
		wantErr: ErrParse,
	}, {
		s:       "cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:",
		wantErr: ErrParse,
	},
	}

	for i, v := range vectors {
		err := ValidateFS(v.s)
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %d, Error: got %v, want %v", i, err, v.wantErr)
		}
	}
}
