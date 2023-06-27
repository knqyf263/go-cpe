package naming

import (
	"fmt"
	"testing"

	"github.com/knqyf263/go-cpe/common"
	"github.com/pkg/errors"
)

func TestUnbindURI(t *testing.T) {
	vectors := []struct {
		s        string
		expected common.WellFormedName
		wantErr  error
	}{{
		s: "cpe:/a:microsoft:internet_explorer%01%01%01%01:?:beta",
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     "microsoft",
			"product":    "internet_explorer????",
			"version":    "?",
			"update":     "beta",
			"edition":    any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
	}, {
		s: "cpe:/a:microsoft:internet_explorer:8.%2a:sp%3f",
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     "microsoft",
			"product":    "internet_explorer",
			"version":    `8\.\*`,
			"update":     `sp\?`,
			"edition":    any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
	}, {
		s: "cpe:/a:microsoft:internet_explorer:8.%02:sp%01",
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     "microsoft",
			"product":    "internet_explorer",
			"version":    `8\.*`,
			"update":     `sp?`,
			"edition":    any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
	}, {
		s: "cpe:/a:hp:insight_diagnostics:7.4.0.1570::~~online~win2003~x64~",
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     "hp",
			"product":    "insight_diagnostics",
			"version":    `7\.4\.0\.1570`,
			"update":     any,
			"edition":    any,
			"sw_edition": "online",
			"target_sw":  "win2003",
			"target_hw":  "x64",
			"other":      any,
			"language":   any,
		},
	}, {
		s: "cpe:/a:%01%01microsoft",
		expected: common.WellFormedName{
			"part":   "a",
			"vendor": "??microsoft",
		},
	}, {
		s:       "cpe:2.3:a:foo\\$bar:insight:7.4.0.1570:-:*:*:online:win2003:x64:*",
		wantErr: common.ErrParse,
	}, {
		s:       "cpe:/a:microsoft:internet_explorer:8.%02:s%01p",
		wantErr: common.ErrParse,
	}, {
		s:       "cpe:/a:micro%02soft",
		wantErr: common.ErrParse,
	}, {
		s: "cpe:/a:micro%01%01:internet%02:%21%22%23%24%25%26%27%28%29%2a%2b%2c:beta-.::online%2f%3a%3b%3c%3d%3e%3f%40%5b%5c%5d%5e",
		expected: common.WellFormedName{
			"part":     "a",
			"vendor":   "micro??",
			"product":  "internet*",
			"version":  `\!\"\#\$\%\&\'\(\)\*\+\,`,
			"update":   `beta\-\.`,
			"language": `online\/\:\;\<\=\>\?\@\[\\\]\^`,
		},
	}, {
		s: "cpe:/a:%60:%7b%7c%7d%7ea",
		expected: common.WellFormedName{
			"part":    "a",
			"vendor":  "\\`",
			"product": `\{\|\}\~a`,
		},
	}, {
		s:       "cpe:/a:hp:insight_diagnostics:7.4.0.1570::~~online~win2003",
		wantErr: common.ErrParse,
	}, {
		s:       "cpe:/a:hp:insight_diagnostics:7.4.0.1570::~~online~win2003~%99~",
		wantErr: common.ErrParse,
	}, {
		s:       `cpe:/a:hp:insight_diagnostics:7.4.0.1570::~~online~win2003~a*b~`,
		wantErr: common.ErrParse,
	}, {
		s:       "cpe:/a:%99",
		wantErr: common.ErrParse,
	}, {
		s:       "cpe:/a:b*c",
		wantErr: common.ErrParse,
	}, {
		s:       "cpe:/z:%01%01microsoft",
		wantErr: common.ErrParse,
	},
	}

	for i, v := range vectors {
		actual, err := UnbindURI(v.s)
		if errors.Cause(err) != v.wantErr {
			t.Errorf("test %d, Error: got %v, want %v", i, errors.Cause(err), v.wantErr)
		}
		if err != nil {
			continue
		}
		if actual.Get("part") != v.expected.Get("part") {
			t.Errorf("test %d, part: got %v, want %v", i, actual.Get("part"), v.expected.Get("part"))
		}
		if actual.Get("vendor") != v.expected.Get("vendor") {
			t.Errorf("test %d, vendor: got %v, want %v", i, actual.Get("vendor"), v.expected.Get("vendor"))
		}
		if actual.Get("product") != v.expected.Get("product") {
			t.Errorf("test %d, product: got %v, want %v", i, actual.Get("product"), v.expected.Get("product"))
		}
		if actual.Get("version") != v.expected.Get("version") {
			t.Errorf("test %d, version: got %v, want %v", i, actual.Get("version"), v.expected.Get("version"))
		}
		if actual.Get("update") != v.expected.Get("update") {
			t.Errorf("test %d, update: got %v, want %v", i, actual.Get("update"), v.expected.Get("update"))
		}
		if actual.Get("edition") != v.expected.Get("edition") {
			t.Errorf("test %d, edition: got %v, want %v", i, actual.Get("edition"), v.expected.Get("edition"))
		}
		if actual.Get("sw_edition") != v.expected.Get("sw_edition") {
			t.Errorf("test %d, sw_edition: got %v, want %v", i, actual.Get("sw_edition"), v.expected.Get("sw_edition"))
		}
		if actual.Get("target_sw") != v.expected.Get("target_sw") {
			t.Errorf("test %d, target_sw: got %v, want %v", i, actual.Get("target_sw"), v.expected.Get("target_sw"))
		}
		if actual.Get("target_hw") != v.expected.Get("target_hw") {
			t.Errorf("test %d, target_hw: got %v, want %v", i, actual.Get("target_hw"), v.expected.Get("target_hw"))
		}
		if actual.Get("other") != v.expected.Get("other") {
			t.Errorf("test %d, other: got %v, want %v", i, actual.Get("other"), v.expected.Get("other"))
		}
		if actual.Get("language") != v.expected.Get("language") {
			t.Errorf("test %d, language: got %v, want %v", i, actual.Get("language"), v.expected.Get("language"))
		}
	}
}

func TestUnbindFS(t *testing.T) {
	vectors := []struct {
		s        string
		expected common.WellFormedName
		wantErr  error
	}{{
		s: "cpe:2.3:a:micr\\?osoft:internet_explorer:8.0.6001:beta:*:*:*:*:*:*",
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     `micr\?osoft`,
			"product":    "internet_explorer",
			"version":    `8\.0\.6001`,
			"update":     "beta",
			"edition":    any,
			"sw_edition": any,
			"target_sw":  any,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
	}, {
		s: `cpe:2.3:a:\$0.99_kindle_books_project:\$0.99_kindle_books:6:*:*:*:*:android:*:*`,
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     `\$0\.99_kindle_books_project`,
			"product":    `\$0\.99_kindle_books`,
			"version":    `6`,
			"update":     any,
			"edition":    any,
			"sw_edition": any,
			"target_sw":  `android`,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
	}, {
		s: `cpe:2.3:a:2glux*:??com_sexypolling??:0.9.1:-:-:*:-:joomla\!:*:*`,
		expected: common.WellFormedName{
			"part":       "a",
			"vendor":     `2glux*`,
			"product":    `??com_sexypolling??`,
			"version":    `0\.9\.1`,
			"update":     na,
			"edition":    na,
			"sw_edition": na,
			"target_sw":  `joomla\!`,
			"target_hw":  any,
			"other":      any,
			"language":   any,
		},
	}, {
		// invalid prefix
		s:       `cpe:2.4:a:2glux:com_sexypolling:0.9.1:-:-:*:-:joomla\!:*:*`,
		wantErr: common.ErrParse,
	}, {
		// embedded unquoted *
		s:       `cpe:2.3:a:2g*lux:com_sexypolling:0.9.1:-:-:*:-:joomla\!:*:*`,
		wantErr: common.ErrParse,
	}, {
		// embedded unquoted ?
		s:       `cpe:2.3:a:2g?lux:com_sexypolling:0.9.1:-:-:*:-:joomla\!:*:*`,
		wantErr: common.ErrParse,
	}, {
		// invalid  part
		s:       `cpe:2.3:z:2glux*:??com_sexypolling??:0.9.1:-:-:*:-:joomla\!:*:*`,
		wantErr: common.ErrParse,
	},
	}

	for i, v := range vectors {
		actual, err := UnbindFS(v.s)
		if errors.Cause(err) != v.wantErr {
			fmt.Println(err)
			t.Errorf("test %d, Error: got %v, want %v", i, errors.Cause(err), v.wantErr)
		}
		if err != nil {
			continue
		}
		if actual.Get("part") != v.expected.Get("part") {
			t.Errorf("test %d, part: got %v, want %v", i, actual.Get("part"), v.expected.Get("part"))
		}
		if actual.Get("vendor") != v.expected.Get("vendor") {
			t.Errorf("test %d, vendor: got %v, want %v", i, actual.Get("vendor"), v.expected.Get("vendor"))
		}
		if actual.Get("product") != v.expected.Get("product") {
			t.Errorf("test %d, product: got %v, want %v", i, actual.Get("product"), v.expected.Get("product"))
		}
		if actual.Get("version") != v.expected.Get("version") {
			t.Errorf("test %d, version: got %v, want %v", i, actual.Get("version"), v.expected.Get("version"))
		}
		if actual.Get("update") != v.expected.Get("update") {
			t.Errorf("test %d, update: got %v, want %v", i, actual.Get("update"), v.expected.Get("update"))
		}
		if actual.Get("edition") != v.expected.Get("edition") {
			t.Errorf("test %d, edition: got %v, want %v", i, actual.Get("edition"), v.expected.Get("edition"))
		}
		if actual.Get("sw_edition") != v.expected.Get("sw_edition") {
			t.Errorf("test %d, sw_edition: got %v, want %v", i, actual.Get("sw_edition"), v.expected.Get("sw_edition"))
		}
		if actual.Get("target_sw") != v.expected.Get("target_sw") {
			t.Errorf("test %d, target_sw: got %v, want %v", i, actual.Get("target_sw"), v.expected.Get("target_sw"))
		}
		if actual.Get("target_hw") != v.expected.Get("target_hw") {
			t.Errorf("test %d, target_hw: got %v, want %v", i, actual.Get("target_hw"), v.expected.Get("target_hw"))
		}
		if actual.Get("other") != v.expected.Get("other") {
			t.Errorf("test %d, other: got %v, want %v", i, actual.Get("other"), v.expected.Get("other"))
		}
		if actual.Get("language") != v.expected.Get("language") {
			t.Errorf("test %d, language: got %v, want %v", i, actual.Get("language"), v.expected.Get("language"))
		}
	}

}

func TestGetCompURI(t *testing.T) {
	vectors := []struct {
		uri      string
		index    int
		expected string
	}{{
		uri:      "cpe:/a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    0,
		expected: "cpe:",
	}, {
		uri:      "cpe:/a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    1,
		expected: "a",
	}, {
		uri:      "cpe:/a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    2,
		expected: "microsoft",
	}, {
		uri:      "cpe:/a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    100,
		expected: "",
	},
	}

	for i, v := range vectors {
		actual := getCompURI(v.uri, v.index)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestGetCompFS(t *testing.T) {
	vectors := []struct {
		fs       string
		index    int
		expected string
	}{{
		fs:       "cpe:2.3:a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    0,
		expected: "cpe",
	}, {
		fs:       "cpe:2.3:a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    1,
		expected: "2.3",
	}, {
		fs:       "cpe:2.3:a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    2,
		expected: "a",
	}, {
		fs:       "cpe:2.3:a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    100,
		expected: "",
	}, {
		fs:       "cpe:2.3:a:microsoft:internet_explorer%01%01%01%01:?:beta",
		index:    -1,
		expected: "",
	}, {
		fs:       `a\:b\:c:d\:e\:f:g`,
		index:    0,
		expected: `a\:b\:c`,
	}, {
		fs:       `a\:b\:c:d\:e\:f:g`,
		index:    1,
		expected: `d\:e\:f`,
	}, {
		fs:       `a\:b\:c:d\:e\:f:g`,
		index:    2,
		expected: "g",
	}, {
		fs:       `a\:b\:c:d\:e\:f:g`,
		index:    3,
		expected: "",
	},
	}

	for i, v := range vectors {
		actual := getCompFS(v.fs, v.index)
		if actual != v.expected {
			t.Errorf("test %d, Result: got %v, want %v", i, actual, v.expected)
		}
	}
}

func TestUnpack(t *testing.T) {
	vectors := []struct {
		s        string
		expected common.WellFormedName
		wantErr  error
	}{{
		s: "~-~-~wordpress~~",
		expected: common.WellFormedName{
			"edition":    na,
			"sw_edition": na,
			"target_sw":  "wordpress",
			"target_hw":  any,
			"other":      any,
		},
	}, {
		s: "~~~~~standalone",
		expected: common.WellFormedName{
			"other": "standalone",
		},
	}, {
		s: "~-~portal~sw~x86~a~b~c",
		expected: common.WellFormedName{
			"edition":    na,
			"sw_edition": "portal",
			"target_sw":  "sw",
			"target_hw":  "x86",
			"other":      `a\~b\~c`,
		},
	},
	}

	for i, v := range vectors {
		wfn := common.WellFormedName{}
		actual, err := unpack(v.s, wfn)
		if errors.Cause(err) != v.wantErr {
			fmt.Println(err)
			t.Errorf("test %d, Error: got %v, want %v", i, errors.Cause(err), v.wantErr)
		}
		if err != nil {
			continue
		}
		if actual.Get("edition") != v.expected.Get("edition") {
			t.Errorf("test %d, edition: got %v, want %v", i, actual.Get("edition"), v.expected.Get("edition"))
		}
		if actual.Get("sw_edition") != v.expected.Get("sw_edition") {
			t.Errorf("test %d, sw_edition: got %v, want %v", i, actual.Get("sw_edition"), v.expected.Get("sw_edition"))
		}
		if actual.Get("target_sw") != v.expected.Get("target_sw") {
			t.Errorf("test %d, target_sw: got %v, want %v", i, actual.Get("target_sw"), v.expected.Get("target_sw"))
		}
		if actual.Get("target_hw") != v.expected.Get("target_hw") {
			t.Errorf("test %d, target_hw: got %v, want %v", i, actual.Get("target_hw"), v.expected.Get("target_hw"))
		}
		if actual.Get("other") != v.expected.Get("other") {
			t.Errorf("test %d, other: got %v, want %v", i, actual.Get("other"), v.expected.Get("other"))
		}
	}

}
