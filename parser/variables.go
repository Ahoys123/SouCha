package parser

import (
	"log"
)

type RuleContext struct {
	Features map[string]*ValueSet
}

// ValueSet is a special type of Set which only accepts Values as elements.
type ValueSet struct {
	set     map[Value]struct{}
	binding string
	not     bool
}

// Union returns a newly allocated map based on the union of both sets.
// The returned set will return values that exist in EITHER set.
func (ms0 *ValueSet) Union(ms1 *ValueSet) *ValueSet {
	if ms0.not == ms1.not && ms0.not {
		// use demorgan's law if both are nots; !1 | !2 = !(1 & 2)
		ms0.not, ms1.not = false, false
		unNot := ms0.Intersection(ms1)
		unNot.not = true
		ms0.not, ms1.not = true, true
		return unNot
	}

	if ms0.not != ms1.not {
		// use property !1 | 2 = !(1 & !2) from demorgans
		var notNot, not *ValueSet
		if ms0.not {
			not, notNot = ms0, ms1
		} else {
			not, notNot = ms1, ms0
		}

		notNot.not, not.not = true, false // flip
		unNot := notNot.Intersection(not)
		unNot.not = true
		notNot.not, not.not = false, true // unflip to previous state
		return unNot
	}

	trms := &ValueSet{set: make(map[Value]struct{})}
	for k := range ms0.set {
		trms.set[k] = struct{}{}
	}
	for k := range ms1.set {
		trms.set[k] = struct{}{}
	}
	return trms
}

// Intersection returns a newly allocated map based on the intersection of both sets.
// The returned set will return values that exist in BOTH sets.
func (ms0 *ValueSet) Intersection(ms1 *ValueSet) *ValueSet {
	if ms0.not == ms1.not && ms0.not {
		// use demorgan's law if both are nots; !1 & !2 = !(1 | 2)
		ms0.not, ms1.not = false, false
		unNot := ms0.Union(ms1)
		unNot.not = true
		ms0.not, ms1.not = true, true
		return unNot
	}

	var notNot, not *ValueSet
	if ms0.not {
		not, notNot = ms0, ms1
	} else {
		not, notNot = ms1, ms0
	}

	trms := &ValueSet{set: make(map[Value]struct{})}
	for k := range notNot.set {
		if _, in := not.set[k]; in != not.not { // if not == true, then if in set == false, add it
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
	var notted bool
	tlen := len(txt)

	// A + !B
	// A ! B
	// -> go that way

	// basic logic: consume characters until an operator is reached;
	// when it is, take the characters between the previous operator or the begining of the string
	// and parse them (idgaf how), then combine them according to the PREVIOUS operator seen
	// then set the previously seen operator to be the current operator.

	lwi := 0 // last written index
	for i := 0; i < tlen; i++ {
		switch txt[i] {
		case '+', '&':
			right = saveRight(txt, i, lwi, right, ctx)
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx, notted)
			notted = false
			right = nil
			oper = intersection
			lwi = i + 1
		case '-', '!': // is equal to + !
			right = saveRight(txt, i, lwi, right, ctx)

			if right != nil { // is something like "Z + A ! B"
				left = operationSwitch(txt, oper, i, lwi, left, right, ctx, notted)
				right = nil
				oper = intersection
			} // else is Z + A + !B

			notted = true
			lwi = i + 1
		case '|':
			right = saveRight(txt, i, lwi, right, ctx)
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx, notted)
			notted = false
			right = nil
			oper = union
			lwi = i + 1
		case ' ':
			right = saveRight(txt, i, lwi, right, ctx)
			lwi = i + 1
		case ':':
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
			right = saveRight(txt, i, lwi, right, ctx)
			return operationSwitch(txt, oper, i, lwi, left, right, ctx, notted), i + 1
		}
	}
	right = saveRight(txt, tlen, lwi, right, ctx)
	left = operationSwitch(txt, oper, tlen, lwi, left, right, ctx, notted)
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
func operationSwitch(txt string, oper uint8, i, lwi int, left *ValueSet, right *ValueSet, ctx *RuleContext, notted bool) *ValueSet {
	if left == nil {
		// CHECK!!!!
		if right == nil {
			return &ValueSet{not: true}
		}
		right.not = notted
		return right
	}

	if right == nil {
		return left
	}

	right.not = notted

	switch oper {
	case union:
		return left.Union(right)
	case intersection:
		return left.Intersection(right)
	}

	log.Fatalln("ERROR!")
	return left
}

// Constant for operations on ValueSets
const (
	union uint8 = iota
	intersection
)

// ValueSetToSet converts a ValueSet to a *Set for use in Matchables.
func ValueSetToSet(ms *ValueSet) *Set {
	s := &Set{[]Matchable{}, ms.binding, ms.not}
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

/*
A{p t k} B{p t k b d g}
A ! B = A + !B
union between !A and B = difference between B and A
*/
