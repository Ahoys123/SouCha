package main

import (
	"main/parser"
)

func main() {
	parser.NewLanguage("[-(stop+labial)]n > m", "pbtdkg")
	// a{b c}d > b

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
