package indexer

import (
	"reflect"
	"sort"
)

type StringSet map[string]struct{}

// https://pkg.go.dev/go.chromium.org/luci/common/data/stringset#Set

// NewStringSet creates a StringSet from a list of values.
func NewStringSet(items ...string) StringSet {
	ss := StringSet{}
	ss.Insert(items...)
	return ss
}

// StringKeySet creates a StringSet from a keys of a map[string](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func StringKeySet(theMap interface{}) StringSet {
	v := reflect.ValueOf(theMap)
	ret := StringSet{}

	for _, keyValue := range v.MapKeys() {
		ret.Insert(keyValue.Interface().(string))
	}
	return ret
}

// Insert adds items to the set.
func (ss StringSet) Insert(items ...string) StringSet {
	for _, item := range items {
		ss[item] = struct{}{}
	}
	return ss
}

// Delete removes all items from the set.
func (ss StringSet) Delete(items ...string) StringSet {
	for _, item := range items {
		delete(ss, item)
	}
	return ss
}

// Has returns true if and only if item is contained in the set.
func (ss StringSet) Has(item string) bool {
	_, contained := ss[item]
	return contained
}

// HasAll returns true if and only if all items are contained in the set.
func (ss StringSet) HasAll(items ...string) bool {
	for _, item := range items {
		if !ss.Has(item) {
			return false
		}
	}
	return true
}

// HasAny returns true if any items are contained in the set.
func (ss StringSet) HasAny(items ...string) bool {
	for _, item := range items {
		if ss.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
func (ss StringSet) Difference(s2 StringSet) StringSet {
	result := NewStringSet()
	for key := range ss {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Union returns a new set which includes items in either s1 or s2.
func (ss StringSet) Union(s2 StringSet) StringSet {
	result := NewStringSet()
	for key := range ss {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
func (ss StringSet) Intersection(s2 StringSet) StringSet {
	var walk, other StringSet
	result := NewStringSet()
	if ss.Len() < s2.Len() {
		walk = ss
		other = s2
	} else {
		walk = s2
		other = ss
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (ss StringSet) IsSuperset(s2 StringSet) bool {
	for item := range s2 {
		if !ss.Has(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (ss StringSet) Equal(s2 StringSet) bool {
	return len(ss) == len(s2) && ss.IsSuperset(s2)
}

type sortableSliceOfString []string

func (s sortableSliceOfString) Len() int           { return len(s) }
func (s sortableSliceOfString) Less(i, j int) bool { return lessString(s[i], s[j]) }
func (s sortableSliceOfString) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// List returns the contents as a sorted string slice.
func (ss StringSet) List() []string {
	res := make(sortableSliceOfString, 0, len(ss))
	for key := range ss {
		res = append(res, key)
	}
	sort.Sort(res)
	return []string(res)
}

// UnsortedList returns the slice with contents in random order.
func (ss StringSet) UnsortedList() []string {
	res := make([]string, 0, len(ss))
	for key := range ss {
		res = append(res, key)
	}
	return res
}

// Returns a single element from the set.
func (ss StringSet) PopAny() (string, bool) {
	for key := range ss {
		ss.Delete(key)
		return key, true
	}
	var zeroValue string
	return zeroValue, false
}

// Len returns the size of the set.
func (ss StringSet) Len() int {
	return len(ss)
}

func lessString(lhs, rhs string) bool {
	return lhs < rhs
}
