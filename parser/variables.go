package parser

import (
	"fmt"
)

type RuleContext struct {
	Features  map[string]MapSet
	Universal MapSet
}

type MapSet map[Value]struct{}

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

// TODO: add user errors
// TODO: add __to__ support
func NewVarSet(txt string, ctx *RuleContext) (MapSet, int) {
	fmt.Println(txt)
	var left, right MapSet
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
			fmt.Println(left)
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
		case '{':
			ms, consumed := NewMapSet(txt[i+1:])
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

	return left, tlen
}

func saveRight(txt string, i, lwi int, right MapSet, ctx *RuleContext) MapSet {
	if i > lwi {
		ms, ok := ctx.Features[txt[lwi:i]]
		if ok {
			return ms
		}
	}
	return right
}

func operationSwitch(txt string, oper uint8, i, lwi int, left MapSet, right MapSet, ctx *RuleContext) MapSet {
	right = saveRight(txt, i, lwi, right, ctx)

	fmt.Printf("\t%s, %s, %d\n", left, right, oper)

	if left == nil {
		if right == nil {
			switch oper {
			case union:
				return ctx.Universal
			case difference:
				return ctx.Universal
			case intersection:
				return MapSet{}
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
