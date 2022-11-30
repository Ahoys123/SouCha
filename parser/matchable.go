package parser

import (
	"fmt"
	"strings"
)

type Matchable interface {
	// MatchStart returns if the begining of text matches.
	MatchStart(txt string) (int, []int, map[string]Value)
	// FollowPath returns the equivielent strucutre
	FollowPath(path []int, bindings map[string]Value) string
}

// NewMatchable creates a new matchable from the string.
func NewMatchable(txt string, ctx *RuleContext) (Matchable, int) {
	set, seq := &Set{}, &Sequence{}
	lci := 0 // last commit index
	for i := 0; i < len(txt); i++ {
		switch txt[i] {
		// if opening bracket, hand off to respective handlers and then continue from after they left off
		case '{', '(', '[':
			savePrev(seq, txt, i, lci)

			var m Matchable
			var last int
			switch txt[i] {
			case '{':
				m, last = NewMatchable(txt[i+1:], ctx)
			case '(':
				m, last = NewOptional(txt[i+1:], ctx)
			case '[':
				last = strings.IndexByte(txt[i+1:], ']') + 1
				vs, _ := NewVarSet(txt[i+1:i+last], ctx)
				m = ValueSetToSet(vs)

			}
			seq.arr = append(seq.arr, collapse(m))

			i += last
			lci = i + 1
		// if closing bracket, return info between start (presumably start of brackets) and end of brackets
		case '}', ')':
			savePrev(seq, txt, i, lci)
			if len(seq.arr) != 0 {
				set.arr = append(set.arr, collapse(seq))
			}
			return set, i + 1
		// if space, save the previous sequence into the return set, then continue with the next sequence
		case ' ', ',':
			savePrev(seq, txt, i, lci)
			if len(seq.arr) != 0 {
				set.arr = append(set.arr, collapse(seq))
				seq = &Sequence{}
			}
			lci = i + 1
		}
	}

	savePrev(seq, txt, len(txt), lci)

	if len(seq.arr) != 0 {
		set.arr = append(set.arr, collapse(seq))
		seq = &Sequence{}
	}

	return collapse(set), len(seq.arr)
}

// NewOptional allows a string of "" to match, along with any others.
func NewOptional(txt string, ctx *RuleContext) (Matchable, int) {
	m, i := NewMatchable(txt, ctx)
	return &Set{[]Matchable{m, Value("")}, ""}, i
}

// savePrev saves the previous value up to that point into seq if that text is not whitespace or empty.
func savePrev(seq *Sequence, txt string, i, lci int) {
	if i > lci {
		seq.arr = append(seq.arr, Value(txt[lci:i]))
	}
}

// collapse turns matchables into reduced forms, with sets or sequences with only one value collapsed to that value.
func collapse(m Matchable) Matchable {
	switch e := m.(type) {
	case *Sequence:
		switch len(e.arr) {
		case 0:
			return nil
		case 1:
			return e.arr[0]
		default:
			return e
		}
	case *Set:
		switch len(e.arr) {
		case 0:
			return nil
		case 1:
			return e.arr[0]
		default:
			return e
		}
	default:
		return e
	}
}

// Sequence satisfies MatchStart when each Matchable in it's array matches, in order, the begining of the text.
type Sequence struct {
	arr []Matchable
}

func (s *Sequence) MatchStart(text string) (int, []int, map[string]Value) {
	i := 0
	totalBindings := map[string]Value{}
	for _, v := range s.arr {
		if consumed, _, bindings := v.MatchStart(text[i:]); consumed != -1 {
			i += consumed
			// add bindings
			for k, v := range bindings {
				totalBindings[k] = v
			}
		} else {
			return -1, nil, nil
		}
	}
	return i, nil, totalBindings
}

func (s *Sequence) FollowPath(path []int, bindings map[string]Value) string {
	return "ERROR"
}

func (s *Sequence) String() string {
	return fmt.Sprintf("Seq:%s", s.arr)
}

// Set satisfies MatchStart when any Matchable in it's array matches the begining of the text.
type Set struct {
	arr     []Matchable
	binding string
}

func (s *Set) MatchStart(text string) (int, []int, map[string]Value) {
	for i, v := range s.arr {
		if consumed, path, bindings := v.MatchStart(text); consumed != -1 {
			if s.binding != "" {
				if bindings == nil {
					bindings = map[string]Value{}
				}
				bindings[s.binding] = Value(text[:consumed])
			}
			return consumed, append(path, i), bindings
		}
	}
	return -1, nil, nil
}

func (s *Set) FollowPath(path []int, bindings map[string]Value) string {
	lasti := len(path) - 1
	fmt.Println(s)

	// if this set has a binding that is defined, return the binding
	if v, ok := bindings[s.binding]; ok {
		return string(v)
	}

	// else, continue following the path
	return s.arr[path[lasti]].FollowPath(path[:lasti], bindings)
}

func (s *Set) String() string {
	return fmt.Sprintf("Set(%s):%s", s.binding, s.arr)
}

// Value satisfies MatchStart when it matches the begining of the text.
type Value string

func (v Value) MatchStart(text string) (int, []int, map[string]Value) {
	vlen, tlen := len(v), len(text)
	vi := 0
	add := 0

	if v == "" {
		return 0, nil, nil
	}

	for i := 0; i < (vlen+add) && i < tlen; i++ {
		switch v[vi] {
		case '#':
			next := waiter(text[i:], ' ')
			if next == -1 {
				return -1, nil, nil
			}
			add += next - 1
			i += next - 1

			vi++
		case text[i]:
			vi++
		default:
			return -1, nil, nil
		}

		if vi >= vlen {
			return vlen + add, nil, nil
		}
	}
	return -1, nil, nil
}

// waiter stalls on a character for however long, returning -1 if the stalled character doesn't exist.
// Otherwise, it returns how many characters it stalled on.
func waiter(txt string, on byte) int {
	if txt[0] != on {
		return -1
	}
	for i := range txt[1:] {
		if txt[1+i] != on {
			return 1 + i
		}
	}
	return len(txt)
}

func (v Value) FollowPath(path []int, bindings map[string]Value) string {
	return string(v)
}

func (v Value) String() string {
	return "'" + string(v) + "'"
}
