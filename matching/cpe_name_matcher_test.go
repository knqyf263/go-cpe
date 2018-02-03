package matching

import (
	"testing"

	"github.com/knqyf263/go-cpe/common"
)

func TestMatcher(t *testing.T) {
	any, _ := common.NewLogicalValue("ANY")
	na, _ := common.NewLogicalValue("NA")
	vectors := []struct {
		wfn                common.WellFormedName
		wfn2               common.WellFormedName
		expectedIsDisjoint bool
		expectedIsEqual    bool
		expectedIsSubset   bool
		expectedIsSuperset bool
	}{{
		wfn: common.WellFormedName{
			"part":       "a",
			"vendor":     `microsoft`,
			"product":    "internet_explorer",
			"version":    `8\.0\.6001`,
			"update":     "beta",
			"edition":    any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   "sp2",
		},
		wfn2: common.WellFormedName{
			"part":       "a",
			"vendor":     `microsoft`,
			"product":    "internet_explorer",
			"version":    any,
			"update":     "beta",
			"edition":    any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   true,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":       "a",
			"vendor":     `adobe`,
			"product":    any,
			"version":    `9\.*`,
			"update":     any,
			"edition":    "PalmOS",
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
		wfn2: common.WellFormedName{
			"part":       "a",
			"vendor":     any,
			"product":    "reader",
			"version":    `9\.3\.2`,
			"update":     na,
			"edition":    na,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":       "a",
			"vendor":     `adobe`,
			"product":    "reader",
			"version":    `9\.*`,
			"update":     na,
			"edition":    na,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
		wfn2: common.WellFormedName{
			"part":       "a",
			"vendor":     `adobe`,
			"product":    "reader",
			"version":    `8\.3\.2`,
			"update":     na,
			"edition":    na,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":       "o",
			"vendor":     `microsoft`,
			"product":    `windows_?`,
			"version":    any,
			"update":     any,
			"edition":    any,
			"sw_edition": `home*`,
			"target_sw":  na,
			"target_hw":  `x64`,
			"other":      na,
			"language":   `en\-us`,
		},
		wfn2: common.WellFormedName{
			"part":       "o",
			"vendor":     `microsoft`,
			"product":    `windows_7`,
			"version":    `6\.1`,
			"update":     "sp1",
			"edition":    any,
			"sw_edition": `home_basic`,
			"target_sw":  na,
			"target_hw":  `x32`,
			"other":      any,
			"language":   `en\-us`,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":       "o",
			"vendor":     `microsoft`,
			"product":    `windows_?`,
			"version":    any,
			"update":     any,
			"edition":    any,
			"sw_edition": `home*`,
			"target_sw":  na,
			"target_hw":  `x64`,
			"other":      na,
			"language":   `en\-us`,
		},
		wfn2: common.WellFormedName{
			"part":       "o",
			"vendor":     `microsoft`,
			"product":    `windows_7`,
			"version":    `6\.1`,
			"update":     "sp1",
			"edition":    any,
			"sw_edition": `home_basic`,
			"target_sw":  na,
			"target_hw":  `x64`,
			"other":      na,
			"language":   `en\-us`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":       "o",
			"vendor":     `microsoft`,
			"product":    `windows_7`,
			"version":    any,
			"update":     any,
			"edition":    any,
			"sw_edition": `home_basic`,
			"target_sw":  na,
			"target_hw":  `x64`,
			"other":      na,
			"language":   `en\-us`,
		},
		wfn2: common.WellFormedName{
			"part":       "o",
			"vendor":     `microsoft`,
			"product":    `windows_7`,
			"version":    any,
			"update":     any,
			"edition":    any,
			"sw_edition": `home_basic`,
			"target_sw":  na,
			"target_hw":  `x64`,
			"other":      na,
			"language":   `en\-us`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    true,
		expectedIsSubset:   true,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": `microsoft`,
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microsoft*`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": na,
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microsoft`,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "*soft",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microsoft`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "?icrosoft",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microsoft`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "?crosoft",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microsoft`,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "??icrosoft",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `\!\#icrosoft`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "microso??",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microso\!\#`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "microso??",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `microso\!\#ft`,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "?",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `a`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "???",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `ab`,
		},
		expectedIsDisjoint: false,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: true,
	}, {
		wfn: common.WellFormedName{
			"part":   "o",
			"vendor": "?",
		},
		wfn2: common.WellFormedName{
			"part":   "o",
			"vendor": `ab`,
		},
		expectedIsDisjoint: true,
		expectedIsEqual:    false,
		expectedIsSubset:   false,
		expectedIsSuperset: false,
	},
	}

	var actual bool
	for i, v := range vectors {
		actual = IsDisjoint(v.wfn, v.wfn2)
		if actual != v.expectedIsDisjoint {
			t.Errorf("test %d, IsDisJoint: got %v, want %v", i, actual, v.expectedIsDisjoint)
		}
		actual = IsEqual(v.wfn, v.wfn2)
		if actual != v.expectedIsEqual {
			t.Errorf("test %d, IsEqual: got %v, want %v", i, actual, v.expectedIsEqual)
		}
		actual = IsSubset(v.wfn, v.wfn2)
		if actual != v.expectedIsSubset {
			t.Errorf("test %d, IsSubset: got %v, want %v", i, actual, v.expectedIsSubset)
		}
		actual = IsSuperset(v.wfn, v.wfn2)
		if actual != v.expectedIsSuperset {
			t.Errorf("test %d, IsSuperset: got %v, want %v", i, actual, v.expectedIsSuperset)
		}
	}
}

func TestIsEvenWildcards(t *testing.T) {
	vectors := []struct {
		str      string
		index    int
		expected bool
	}{{
		str:      `abc`,
		index:    2,
		expected: true,
	}, {
		str:      `abc*`,
		index:    3,
		expected: true,
	}, {
		str:      `abc\*`,
		index:    4,
		expected: false,
	}, {
		str:      `abc\\*`,
		index:    5,
		expected: true,
	}, {
		str:      `abc\\\*`,
		index:    6,
		expected: false,
	}, {
		str:      `abc\\def\*`,
		index:    9,
		expected: false,
	},
	}

	for i, v := range vectors {
		actual := IsEvenWildcards(v.str, v.index)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}
