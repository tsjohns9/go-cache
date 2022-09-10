package cache

type Operator string

const (
	DoubleEquals Operator = "=="
	Exists       Operator = "exists"
)

type AttrMap map[string]string

type Tags []string

// Selector represents a label selector.
type Selector interface {
	Matches(Items) bool
	Empty() bool
	Add(r ...Rule) Selector
}

type Rule interface {
	Matches(items Items) bool
}

// Items -
type Items interface {
	// Has returns whether the provided value exists as an attribute or tag
	Has(v string) (exists bool)
	// Get returns the value for the provided attribute, or the tag.
	Get(v string) (value string)
	Kind() string
}

type selector struct {
	rules []Rule
}

func NewSelector() Selector {
	return &selector{}
}

func (s *selector) Matches(items Items) bool {
	for _, rule := range s.rules {
		if !rule.Matches(items) {
			return false
		}
	}
	return true
}

func (s *selector) Empty() bool {
	return len(s.rules) == 0
}

func (s *selector) Add(r ...Rule) Selector {
	s.rules = append(s.rules, r...)
	return s
}

type AttrRule struct {
	key   string
	value []string
	op    Operator
}

type TagRule struct {
	value string
	op    Operator
}

func (r *AttrRule) Matches(items Items) bool {
	if items.Kind() != "attributes" {
		return false
	}

	switch r.op {
	case Exists:
		return items.Has(r.key)

	case DoubleEquals:
		v := items.Get(r.key)

		for _, val := range r.value {
			if v == val {
				return true
			}
		}

		return false
	default:
		return false
	}
}

func NewAttrRule(key string, op Operator, value []string) *AttrRule {
	return &AttrRule{
		key: key, value: value, op: op,
	}
}

func NewTagRule(value string) *TagRule {
	return &TagRule{
		value: value,
		op:    DoubleEquals,
	}
}

func (r *TagRule) Matches(items Items) bool {
	if items.Kind() != "tags" {
		return false
	}

	switch r.op {
	case DoubleEquals:
		return items.Has(r.value)
	default:
		return false
	}
}

// Has returns whether the provided label exists in the map.
func (t Tags) Has(label string) bool {
	for _, v := range t {
		if v == label {
			return true
		}
	}
	return false
}

// Get returns the value in the map for the provided label.
func (t Tags) Get(val string) string {
	for _, v := range t {
		if v == val {
			return v
		}
	}
	return ""
}

func (t Tags) Kind() string {
	return "tags"
}

// Has returns whether the provided label exists in the map.
func (a AttrMap) Has(label string) bool {
	_, exists := a[label]
	return exists
}

// Get returns the value in the map for the provided label.
func (a AttrMap) Get(label string) string {
	return a[label]
}

func (a AttrMap) Kind() string {
	return "attributes"
}
