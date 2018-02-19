package matching

import (
	"strings"

	"github.com/knqyf263/go-cpe/common"
)

// Relation is enumeration for relational values.
type Relation int

const (
	// DISJOINT : disjoint
	DISJOINT Relation = iota
	// SUBSET : subset
	SUBSET
	// SUPERSET : superset
	SUPERSET
	// EQUAL : equal
	EQUAL
	// UNDEFINED : undefined
	UNDEFINED
)

// IsDisjoint tests two Well Formed Names for disjointness.
// @param source Source WFN
// @param target Target WFN
// @return true if the names are disjoint, false otherwise
func IsDisjoint(source, target common.WellFormedName) bool {
	// if any pairwise comparison is disjoint, the names are disjoint.
	results := CompareWFNs(source, target)
	for _, result := range results {
		if result == DISJOINT {
			return true
		}
	}
	return false
}

// IsEqual tests two Well Formed Names for equality.
// @param source Source WFN
// @param target Target WFN
// @return true if the names are equal, false otherwise
func IsEqual(source, target common.WellFormedName) bool {
	// if every pairwise comparison is equal, the names are equal.
	results := CompareWFNs(source, target)
	for _, result := range results {
		if result != EQUAL {
			return false
		}
	}
	return true
}

// IsSubset tests if the source Well Formed Name is a subset of the target Well Formed
// Name.
// @param source Source WFN
// @param target Target WFN
// @return true if the source is a subset of the target, false otherwise
func IsSubset(source, target common.WellFormedName) bool {
	// if any comparison is anything other than subset or equal, then target is
	// not a subset of source.
	results := CompareWFNs(source, target)
	for _, result := range results {
		if result != SUBSET && result != EQUAL {
			return false
		}
	}
	return true
}

// IsSuperset tests if the source Well Formed name is a superset of the target Well Formed Name.
// @param source Source WFN
// @param target Target WFN
// @return true if the source is a superset of the target, false otherwise
func IsSuperset(source, target common.WellFormedName) bool {
	// if any comparison is anything other than superset or equal, then target is not
	// a superset of source.
	results := CompareWFNs(source, target)
	for _, result := range results {
		if result != SUPERSET && result != EQUAL {
			return false
		}
	}
	return true
}

// CompareWFNs compares each attribute value pair in two Well Formed Names.
// @param source Source WFN
// @param target Target WFN
// @return A Hashtable mapping attribute string to attribute value Relation
func CompareWFNs(source, target common.WellFormedName) map[string]Relation {
	result := map[string]Relation{}
	result[common.AttributePart] = compare(source.Get(common.AttributePart), target.Get(common.AttributePart))
	result[common.AttributeVendor] = compare(source.Get(common.AttributeVendor), target.Get(common.AttributeVendor))
	result[common.AttributeProduct] = compare(source.Get(common.AttributeProduct), target.Get(common.AttributeProduct))
	result[common.AttributeVersion] = compare(source.Get(common.AttributeVersion), target.Get(common.AttributeVersion))
	result[common.AttributeUpdate] = compare(source.Get(common.AttributeUpdate), target.Get(common.AttributeUpdate))
	result[common.AttributeEdition] = compare(source.Get(common.AttributeEdition), target.Get(common.AttributeEdition))
	result[common.AttributeLanguage] = compare(source.Get(common.AttributeLanguage), target.Get(common.AttributeLanguage))
	result[common.AttributeSwEdition] = compare(source.Get(common.AttributeSwEdition), target.Get(common.AttributeSwEdition))
	result[common.AttributeTargetSw] = compare(source.Get(common.AttributeTargetSw), target.Get(common.AttributeTargetSw))
	result[common.AttributeTargetHw] = compare(source.Get(common.AttributeTargetHw), target.Get(common.AttributeTargetHw))
	result[common.AttributeOther] = compare(source.Get(common.AttributeOther), target.Get(common.AttributeOther))
	return result
}

// Compares an attribute value pair.
// @param source Source attribute value.
// @param target Target attribute value.
// @return The relation between the two attribute values.
func compare(source, target interface{}) Relation {
	var s, t string
	var ok bool

	// matching is case insensitive, convert strings to lowercase.
	if s, ok = source.(string); ok {
		s = strings.ToLower(s)
	}
	if t, ok = target.(string); ok {
		t = strings.ToLower(t)
	}
	// Unquoted wildcard characters yield an undefined result.
	if common.ContainsWildcards(t) {
		return UNDEFINED
	}

	// If source and target values are equal, then result is equal.
	if source == target {
		return EQUAL
	}

	// Check to see if source or target are Logical Values.
	var lvSource, lvTarget common.LogicalValue
	if lv, ok := source.(common.LogicalValue); ok {
		lvSource = lv
	}
	if lv, ok := target.(common.LogicalValue); ok {
		lvTarget = lv
	}
	// If source value is ANY, result is a superset.
	if lvSource.IsANY() {
		return SUPERSET
	}
	// If target value is ANY, result is a subset.
	if lvTarget.IsANY() {
		return SUBSET
	}
	// If source or target is NA, result is disjoint.
	if lvSource.IsNA() {
		return DISJOINT
	}
	// if (lvTarget != null) {
	if lvTarget.IsNA() {
		return DISJOINT
	}
	// only Strings will get to this point, not LogicalValues
	return compareStrings(s, t)
}

// compareStrings compares a source string to a target string, and addresses the condition
// in which the source string includes unquoted special characters. It
// performs a simple regular expression  match, with the assumption that
// (as required) unquoted special characters appear only at the beginning
// and/or the end of the source string. It also properly differentiates
// between unquoted and quoted special characters.
//
// @return Relation between source and target Strings.
func compareStrings(source, target string) Relation {
	var start, begins, ends, index, leftover, escapes int
	end := len(source)

	if source[0] == '*' {
		start = 1
		begins = -1
	} else {
		for start < len(source) && source[start] == '?' {
			start++
			begins++
		}
	}

	if source[end-1] == '*' && IsEvenWildcards(source, end-1) {
		end--
		ends = -1
	} else {
		for end > 0 && source[end-1] == '?' && IsEvenWildcards(source, end-1) {
			end--
			ends++
		}
	}

	// only ? (e.g. "???")
	if strings.Trim(source, "?") == "" {
		if len(source) >= common.LengthWithEscapeCharacters(target) {
			return SUPERSET
		}
		return DISJOINT
	}
	source = source[start:end]
	index = -1
	leftover = len(target)
	for leftover > 0 {
		index = common.IndexOf(target, source, index+1)
		if index == -1 {
			break
		}
		escapes = common.CountEscapeCharacters(target[:index])
		if index > 0 && begins != -1 && begins < (index-escapes) {
			break
		}
		escapes = common.CountEscapeCharacters(target[index+1:])
		leftover = len(target) - index - escapes - len(source)
		if leftover > 0 && ends != -1 && leftover > ends {
			continue
		}
		return SUPERSET
	}
	return DISJOINT
}

// IsEvenWildcards searches a string for the backslash character
// @param str string to search in
// @param idx end index
// @return true if the number of backslash characters is even, false if odd
func IsEvenWildcards(str string, idx int) bool {
	result := 0
	for idx > 0 && str[idx-1] == '\\' {
		idx--
		result++
	}
	return result%2 == 0
}
