package parser

type Node interface {
	GetValue() []Value
}

type Set []Node

func (s Set) GetValue() []Value {
	tr := []Value{}
	for _, c := range s {
		tr = append(tr, c.GetValue()...)
	}
	return tr
}

type Value string

func (v Value) String() string {
	return "\"" + string(v) + "\""
}

func (v Value) GetValue() []Value {
	return []Value{v}
}

func setify(x string) (Node, int) {
	cons := Set{}
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
				cons = append(cons, Value(x[start:i]))
			}
			return reduce(cons), i
		case ' ', ',':
			if i > start {
				cons = append(cons, Value(x[start:i]))
			}
			start = i + 1
		}
	}

	if len(x) > start+1 {
		cons = append(cons, Value(x[start:]))
	}

	return reduce(cons), -1
}

func reduce(s Set) Node {
	switch len(s) {
	case 0:
		return nil
	case 1:
		return s[0]
	}
	return s
}
