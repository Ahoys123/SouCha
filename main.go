package main

import (
	"fmt"
	"main/parser"
)

func main() {
	l := parser.NewLanguage(`stop = p t k b d g
	consonant = p t k b d g m n j
	vowel = a e i o u
	alveolar = t d n
	t > t /
	a > / r _ r`)
	// a{b c}d > b
	fmt.Println(l.Evolve("atsʼari"))
}

// [+stop+consonant+alveolar] > r / [+vowel+stress] _ [+vowel-stress]
// t d > r / [+vowel+stress] _ [+vowel-stress]

/* plan for world variable domination(?!):


PLAN 1:
if [] in "to":
	h0 = hash corresponding match in "from"
	for each tag:
		switch tagOperator:
			case -:
				set h0 flag "tag" to 0
			case +:
				set h0 flag "tag" to 1
	find closest "real hash" to h0
	replace [] with that

PLAN 2:
predefined:
	groups of tags; "@place bilabial = p b m"; "@manner stop = p b"; "voice = b d g m n ng"
if [] in "to":
	for each tag:
		switch tagOperator:
			case +:
				check tag's group;
*/
