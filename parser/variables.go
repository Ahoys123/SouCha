package parser

import "strings"

func NewVarSet(txt string) (*Set, int) {
	return &Set{[]Matchable{&Value{"hello"}}}, strings.IndexByte(txt, ']') + 1
}
