package parser

import (
	"fmt"
	"strings"
)

type RuleContext struct {
	Features  map[string]MapSet
	Universal MapSet
}

type MapSet map[Value]struct{}

func (ms MapSet) Evaluate() MapSet {
	return ms
}

func (ms0 MapSet) Union(ms1 MapSet) MapSet {
	trms := make(MapSet)
	for k := range ms0 {
		trms[k] = struct{}{}
	}
	for k := range ms1 {
		trms[k] = struct{}{}
	}
	return trms
}

func (ms0 MapSet) Difference(ms1 MapSet) MapSet {
	trms := make(MapSet)
	for k := range ms0 {
		if _, in := ms1[k]; !in {
			trms[k] = struct{}{}
		}
	}
	return trms
}

func (ms0 MapSet) Intersection(ms1 MapSet) MapSet {
	trms := make(MapSet)
	for k := range ms0 {
		if _, in := ms1[k]; in {
			trms[k] = struct{}{}
		}
	}
	return trms
}

func (ctx RuleContext) NewVarSet(txt string) MapSet {
	var left MapSet
	lci := 0
	tlen := len(txt)

	for i := 0; i < tlen; i++ {
		switch txt[i] {
		case '{':
			l, consumed := NewMapSet(txt[i+1:])
			left = l
			i += consumed - 1
		case '!', '&', '|', '+', '-':
			if left == nil {
				left = ctx.Features[strings.TrimSpace(txt[lci:i])]
			}

			fmt.Println(txt[lci:i], txt[i+1:])
			switch txt[i] {
			case '&', '+':
				if len(left) == 0 {
					left = nil
					continue
				}
				return left.Intersection(ctx.NewVarSet(txt[i+1:])) // left + right
			case '|':
				if len(left) == 0 {
					left = nil
					continue
				}
				return left.Union(ctx.NewVarSet(txt[i+1:])) // left U right
			case '!', '-':
				if len(left) == 0 {
					return ctx.Universal.Difference(ctx.NewVarSet(txt[i+1:]))
				}
				return left.Difference(ctx.NewVarSet(txt[i+1:]))
			}
			fmt.Printf("ERROR, unknown operator %s\n", txt[i:i+1])
			return left
		case '(':
			left = ctx.NewVarSet(txt[i+1:])
		case ')':
			tlen = i
		}
	}

	if left == nil && tlen > lci {
		left = ctx.Features[strings.TrimSpace(txt[lci:tlen])]
	}

	return left
}

func MapSetToSet(ms MapSet) *Set {
	s := &Set{[]Matchable{}}
	for k := range ms {
		s.arr = append(s.arr, k)
	}
	return s
}

func NewMapSet(txt string) (MapSet, int) {
	trms := make(MapSet)
	lci := 0
	for i := 0; i < len(txt); i++ {
		switch txt[i] {
		case ' ', ',':
			if i > lci {
				trms[Value(txt[lci:i])] = struct{}{}
			}

			lci = i + 1
		case '}':
			if i > lci {
				trms[Value(txt[lci:i])] = struct{}{}
			}
			return trms, i + 1
		}
	}
	fmt.Println("uh, no closing }")
	return trms, len(txt)
}
