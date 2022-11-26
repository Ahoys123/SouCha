package parser

type RuleContext struct {
	Features  map[string]*ValueSet
	Universal *ValueSet
}

// ValueSet is a special type of Set which only accepts Values as elements.
type ValueSet struct {
	set     map[Value]struct{}
	binding string
}

// Union returns a newly allocated map based on the union of both sets.
// The returned set will return values that exist in EITHER set.
func (ms0 *ValueSet) Union(ms1 *ValueSet) *ValueSet {
	trms := &ValueSet{set: make(map[Value]struct{})}
	for k := range ms0.set {
		trms.set[k] = struct{}{}
	}
	for k := range ms1.set {
		trms.set[k] = struct{}{}
	}
	return trms
}

// Difference returns a newly allocated map based on the difference of both sets.
// The returned set will return values that exist in the first set, but NOT the second one.
func (ms0 *ValueSet) Difference(ms1 *ValueSet) *ValueSet {
	trms := &ValueSet{set: make(map[Value]struct{})}
	for k := range ms0.set {
		if _, in := ms1.set[k]; !in {
			trms.set[k] = struct{}{}
		}
	}
	return trms
}

// Intersection returns a newly allocated map based on the intersection of both sets.
// The returned set will return values that exist in BOTH sets.
func (ms0 *ValueSet) Intersection(ms1 *ValueSet) *ValueSet {
	trms := &ValueSet{set: make(map[Value]struct{})}
	for k := range ms0.set {
		if _, in := ms1.set[k]; in {
			trms.set[k] = struct{}{}
		}
	}
	return trms
}

// NewVarSet creates a ValueSet based on the user defined sets defined in ctx.
// It supports operations like the Union, Difference, and Intersection of user defined sets.
func NewVarSet(txt string, ctx *RuleContext) (*ValueSet, int) {
	var left, right *ValueSet
	var binding string
	var oper uint8
	tlen := len(txt)

	lwi := 0 // last written index
	for i := 0; i < tlen; i++ {
		switch txt[i] {
		case '+', '&':
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx)
			right = nil
			oper = intersection
			lwi = i + 1
		case '-', '!':
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx)
			//fmt.Println(left)
			right = nil
			oper = difference
			lwi = i + 1
		case '|':
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx)
			right = nil
			oper = union
			lwi = i + 1
		case ' ':
			right = saveRight(txt, i, lwi, right, ctx)
			lwi = i + 1
		case '@':
			if i > lwi {
				binding = txt[lwi:i]
			}
			lwi = i + 1
		case '{':
			ms, consumed := NewValueSet(txt[i+1:])
			right = ms
			i += consumed - 1
		case '(':
			ms, consumed := NewVarSet(txt[i+1:], ctx)
			right = ms
			i += consumed
		case ')':
			return operationSwitch(txt, oper, i, lwi, left, right, ctx), i + 1
		}
	}

	left = operationSwitch(txt, oper, tlen, lwi, left, right, ctx)
	left.binding = binding

	return left, tlen
}

// saveRight returns the values to the right of i if "worth saving", otherwise it returns the previous right value.
// Values are "worth saving" if they exist as a user defined set and are not whitespace.
func saveRight(txt string, i, lwi int, right *ValueSet, ctx *RuleContext) *ValueSet {
	if i > lwi {
		ms, ok := ctx.Features[txt[lwi:i]]
		if ok {
			return ms
		}
	}
	return right
}

// operationSwitch consolodates left and right into one set through the operation defined by oper.
// It returns what the new left should be, as well as making right safe to clear.
func operationSwitch(txt string, oper uint8, i, lwi int, left *ValueSet, right *ValueSet, ctx *RuleContext) *ValueSet {
	right = saveRight(txt, i, lwi, right, ctx)

	//fmt.Printf("\t%s, %s, %d\n", left, right, oper)

	if left == nil {
		if right == nil {
			switch oper {
			case union:
				return ctx.Universal
			case difference:
				return ctx.Universal
			case intersection:
				return &ValueSet{}
			}
		}
		return right
	}

	if right == nil {
		return left
	}

	switch oper {
	case union:
		return left.Union(right)
	case difference:
		return left.Difference(right)
	case intersection:
		return left.Intersection(right)
	}

	return left
}

// Constant for operations on ValueSets
const (
	union uint8 = iota
	difference
	intersection
)

// ValueSetToSet converts a ValueSet to a *Set for use in Matchables.
func ValueSetToSet(ms *ValueSet) *Set {
	s := &Set{[]Matchable{}, ms.binding}
	for k := range ms.set {
		s.arr = append(s.arr, k)
	}
	return s
}

// NewValueSet creates a ValueSet out of PURE values, no variables (use NewVarSet for that).
// All comma/space seperated strings will be counted literally.
func NewValueSet(txt string) (*ValueSet, int) {
	trms := &ValueSet{set: make(map[Value]struct{})}
	lci := 0
	for i := 0; i < len(txt); i++ {
		switch txt[i] {
		case ' ', ',':
			if i > lci {
				trms.set[Value(txt[lci:i])] = struct{}{}
			}

			lci = i + 1
		case '}':
			if i > lci {
				trms.set[Value(txt[lci:i])] = struct{}{}
			}
			return trms, i + 1
		}
	}
	if len(txt) > lci {
		trms.set[Value(txt[lci:])] = struct{}{}
	}
	return trms, len(txt)
}
