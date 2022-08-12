package parser

import (
	"fmt"
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

func NewVarSet(txt string, ctx *RuleContext) MapSet {
	var left MapSet
	var right string
	var oper uint8
	lwi := 0 // last written index
	for i := 0; i < len(txt); i++ {
		switch txt[i] {
		case '+', '&':
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx, ctx.Universal)
			oper = intersection
			lwi = i + 1
		case '-', '!':
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx, ctx.Universal)
			oper = difference
			lwi = i + 1
		case '|':
			left = operationSwitch(txt, oper, i, lwi, left, right, ctx, MapSet{})
			oper = union
			lwi = i + 1
		case ' ':
			if i > lwi {
				right = txt[lwi:i]
			}
			lwi = i + 1
		}
	}

	if len(txt) > lwi {
		right = txt[lwi:]
	}

	left = operationSwitch(txt, oper, len(txt), lwi, left, right, ctx, MapSet{})

	return left
}

func operationSwitch(txt string, oper uint8, i, lwi int, left MapSet, right string, ctx *RuleContext, ielwi MapSet) MapSet {
	if right == "" {
		return ielwi
	}

	if i > lwi {
		right = txt[lwi:i]
	}

	fmt.Println(right)

	set := ctx.Features[right]

	if left == nil {
		return set
	}

	switch oper {
	case union:
		return left.Union(set)
	case difference:
		return left.Difference(set)
	case intersection:
		return left.Intersection(set)
	}

	return left

}

const (
	union uint8 = iota
	difference
	intersection
)

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
