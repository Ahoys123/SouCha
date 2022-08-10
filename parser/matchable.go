package parser

import (
	"fmt"
)

type Matchable interface {
	// MatchStart returns if the begining of text matches.
	MatchStart(txt string) (int, []int)
	// FollowPath returns the equivielent strucutre
	FollowPath(path []int) string
}

// NewMatchable creates a new matchable
func NewMatchable(txt string) (Matchable, int) {
	set, seq := &Set{}, &Sequence{}
	lci := 0 // last commit index
	for i := 0; i < len(txt); i++ {
		switch txt[i] {
		case '{':
			savePrev(seq, txt, i, lci)

			m, last := NewMatchable(txt[i+1:])
			seq.arr = append(seq.arr, collapse(m))

			i += last
			lci = i + 1
		case '}', ')':
			savePrev(seq, txt, i, lci)
			if len(seq.arr) != 0 {
				set.arr = append(set.arr, collapse(seq))
			}
			return set, i + 1
		case '(':
			savePrev(seq, txt, i, lci)
			o, last := NewOptional(txt[i+1:])
			seq.arr = append(seq.arr, collapse(o))
			i += last
			lci = i + 1
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

func NewOptional(txt string) (Matchable, int) {
	m, i := NewMatchable(txt)
	return &Set{[]Matchable{m, &Value{""}}}, i
}

func savePrev(seq *Sequence, txt string, i, lci int) {
	if i > lci {
		seq.arr = append(seq.arr, &Value{txt[lci:i]})
	}
}

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

// To satisfy MatchStart, must match each element sucessively
type Sequence struct {
	arr []Matchable
}

func (s *Sequence) MatchStart(text string) (int, []int) {
	i := 0
	for _, v := range s.arr {
		if consumed, _ := v.MatchStart(text[i:]); consumed != -1 {
			i += consumed
		} else {
			return -1, nil
		}
	}
	return i, nil
}

func (s *Sequence) FollowPath(path []int) string {
	return "ERROR"
}

func (s *Sequence) String() string {
	return fmt.Sprintf("Seq:%s", s.arr)
}

// To satisfy MatchStart, must have at least one element in set
type Set struct {
	arr []Matchable
}

func (s *Set) MatchStart(text string) (int, []int) {
	for i, v := range s.arr {
		if consumed, path := v.MatchStart(text); consumed != -1 {
			return consumed, append(path, i)
		}
	}
	return -1, nil
}

func (s *Set) FollowPath(path []int) string {
	lasti := len(path) - 1
	return s.arr[path[lasti]].FollowPath(path[:lasti])
}

func (s *Set) String() string {
	return fmt.Sprintf("Set:%s", s.arr)
}

// To satsify MatchStart, must be value
type Value struct {
	v string
}

func (v *Value) MatchStart(text string) (int, []int) {
	vlen, tlen := len(v.v), len(text)
	vi := 0
	add := 0

	if v.v == "" {
		return 0, nil
	}

	for i := 0; i < (vlen+add) && i < tlen; i++ {
		switch v.v[vi] {
		case '#':
			next := waiter(text[i:], ' ')
			if next == -1 {
				return -1, nil
			}
			add += next - 1
			i += next - 1

			vi++
		case text[i]:
			vi++
		default:
			return -1, nil
		}

		if vi >= vlen {
			return vlen + add, nil
		}
	}
	return -1, nil
}

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

func (v *Value) FollowPath(path []int) string {
	return v.v
}

func (v *Value) String() string {
	return "'" + v.v + "'"
}
