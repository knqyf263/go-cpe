package naming

import (
	"testing"

	"github.com/knqyf263/go-cpe/common"
)

var (
	any, _ = common.NewLogicalValue("ANY")
	na, _  = common.NewLogicalValue("NA")
)

func TestBindToURI(t *testing.T) {
	vectors := []struct {
		w        common.WellFormedName
		expected string
	}{{
		w: common.WellFormedName{
			"part": "a",
		},
		expected: "cpe:/a",
	}, {
		w: common.WellFormedName{
			"part":     "a",
			"vendor":   "microsoft",
			"product":  "internet_explorer",
			"version":  "8\\.0\\.6001",
			"update":   "beta",
			"edition":  any,
			"language": "sp2",
		},
		expected: "cpe:/a:microsoft:internet_explorer:8.0.6001:beta::sp2",
	}, {
		w: common.WellFormedName{
			"part":       "a",
			"vendor":     "foo\\$bar",
			"product":    "insight",
			"version":    "7\\.4\\.0\\.1570",
			"update":     na,
			"sw_edition": "online",
			"target_sw":  "win2003",
			"target_hw":  "x64",
		},
		expected: "cpe:/a:foo%24bar:insight:7.4.0.1570:-:~~online~win2003~x64~",
	}, {
		w: common.WellFormedName{
			"part":   "a",
			"vendor": 1,
		},
		expected: "cpe:/a",
	}, {
		w: common.WellFormedName{
			"part":     "a",
			"vendor":   "micro??",
			"product":  "internet*",
			"version":  `\!\"\#\$\%\&\'\(\)\*\+\,`,
			"update":   `beta\-\.`,
			"language": `online\/\:\;\<\=\>\?\@\[\\\]\^`,
		},
		expected: "cpe:/a:micro%01%01:internet%02:%21%22%23%24%25%26%27%28%29%2a%2b%2c:beta-.::online%2f%3a%3b%3c%3d%3e%3f%40%5b%5c%5d%5e",
	}, {
		w: common.WellFormedName{
			"part":    "a",
			"vendor":  "\\`",
			"product": `\{\|\}\~\a`,
		},
		expected: "cpe:/a:%60:%7b%7c%7d%7ea",
	},
	}

	for i, v := range vectors {
		actual := BindToURI(v.w)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestBindToFS(t *testing.T) {
	vectors := []struct {
		w        common.WellFormedName
		expected string
	}{{
		w: common.WellFormedName{
			"part": "a",
		},
		expected: "cpe:2.3:a:*:*:*:*:*:*:*:*:*:*",
	}, {
		w: common.WellFormedName{
			"part":     "a",
			"vendor":   "microsoft",
			"product":  "internet_explorer",
			"version":  "8\\.0\\.6001",
			"update":   "beta",
			"edition":  any,
			"language": "sp2",
		},
		expected: "cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:sp2:*:*:*:*",
	}, {
		w: common.WellFormedName{
			"part":       "a",
			"vendor":     "foo\\$bar",
			"product":    "insight",
			"version":    "7\\.4\\.0\\.1570",
			"update":     na,
			"sw_edition": "online",
			"target_sw":  "win2003",
			"target_hw":  "x64",
		},
		expected: "cpe:2.3:a:foo\\$bar:insight:7.4.0.1570:-:*:*:online:win2003:x64:*",
	}, {
		w: common.WellFormedName{
			"part":   "a",
			"vendor": 1,
		},
		expected: "cpe:2.3:a::*:*:*:*:*:*:*:*:*",
	},
	}

	for i, v := range vectors {
		actual := BindToFS(v.w)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}
