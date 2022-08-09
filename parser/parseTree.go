package parser

type Sequence []Tree

func (s Sequence) GetValue() []value {
	tr := []value{}
	for _, c := range s {
		tr = append(tr, c.GetValue()...)
	}
	return tr
}

type Tree interface {
	GetValue() []value
}

type set []Tree

func (s set) GetValue() []value {
	tr := []value{}
	for _, c := range s {
		tr = append(tr, c.GetValue()...)
	}
	return tr
}

type value string

func (v value) GetValue() []value {
	return []value{v}
}

func MatchStart(s Tree, text string) ([]int, int) {
	switch e := s.(type) {
	case value:
		elen := len(e)
		if len(text) >= elen && text[:elen] == string(e) {
			return []int{}, elen
		}
		return nil, -1
	case set:
		for i, v := range e {
			path, mlen := MatchStart(v, text)
			if mlen != -1 {
				return append(path, i), mlen
			}
		}
		return nil, -1
	case nil:
		return nil, 0
	}
	return nil, -1
}

func FollowPath(s Tree, path []int) value {
	switch e := s.(type) {
	case set:
		return FollowPath(e[path[len(path)-1]], path[:len(path)-1])
	case value:
		return e
	}
	return ""
}

// setify parses a string and returns a tree and the number of characters it consumed from the input.
//
// returns -1 if entire input was consumed
func setify(x string) (Tree, int) {
	cons := set{}
	start := 0
	for i := 0; i < len(x); i++ {
		switch x[i] {
		case '{':
			n, last := setify(x[i+1:])
			if n != nil {
				cons = append(cons, n)
			}
			i += last + 1
			start = i + 1

		case '}':
			if i > start {
				cons = append(cons, value(x[start:i]))
			}
			return reduce(cons), i
		case ' ', ',':
			if i > start {
				cons = append(cons, value(x[start:i]))
			}
			start = i + 1
		}
	}

	if len(x) > start {
		cons = append(cons, value(x[start:]))
	}

	return reduce(cons), -1
}

func reduce(s set) Tree {
	switch len(s) {
	case 0:
		return nil
	case 1:
		return s[0]
	}
	return s
}
