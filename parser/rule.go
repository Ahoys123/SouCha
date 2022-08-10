package parser

import (
	"fmt"
	"strings"
)

type Rule struct {
	from, to          Matchable
	env               string
	precond, postcond Matchable
}

func NewRule(rule string) *Rule {
	r := &Rule{}

	r.split(rule)

	r.parseEnv()

	return r
}

func (r *Rule) Apply(text string) string {
	b := strings.Builder{}

	stages := [3]Matchable{r.precond, r.from, r.postcond}
	cstage := 0

	var path []int
	p0, p1 := 0, 0
	lastWritten := 0

	// while in bounds of text
	for i := 0; i < len(text); {

	stageInc:
		mlen := 0
		var p []int
		if stages[cstage] != nil {
			mlen, p = stages[cstage].MatchStart(text[i:])
		}

		//fmt.Printf("%s; %d; %d\n", text[i:], cstage, mlen)

		if mlen != -1 {
			//fmt.Printf("%s; %d; %d\n", text[i:], cstage, mlen)
			cstage++

			i += mlen

			//fmt.Println(text[p1:])

			switch cstage {
			case 1:
				p0 = i
			case 2:
				p1 = i
				path = p
			case 3:
				cstage = 0
				b.WriteString(text[lastWritten:p0])
				if r.to != nil {
					b.WriteString(r.to.FollowPath(path))
				}
				i = p1
				lastWritten = p1
			}

			if stages[cstage] == nil {
				goto stageInc
			}

		} else {

			// if failed on postcondition, we need to recheck the condition; else, move on to next
			if cstage != 2 {
				i++
			}
			cstage = 0
		}
	}

	b.WriteString(text[lastWritten:])

	return b.String()
}

func (r *Rule) split(rule string) {
	pointer := 0
	var sto string
	for i := 0; i < len(rule); i++ {
		switch rule[i] {
		case '>':
			if r.from == nil {
				r.from, _ = NewMatchable(strings.TrimSpace(rule[pointer:i]))
			}
			pointer = i + 1
		case '/':
			sto = rule[pointer:i]
			pointer = i + 1
		}
	}

	if sto == "" {
		sto = rule[pointer:]
		r.env = "_"
	} else {
		r.env = strings.TrimSpace(rule[pointer:])
	}

	r.to, _ = NewMatchable(strings.TrimSpace(sto))
}

func (r *Rule) parseEnv() error {
	split := strings.SplitN(r.env, "_", 2)
	if len(split) != 2 {
		return fmt.Errorf("enviornment \"%s\" not in format \"{precondition} _ {postcondition}\"", r.env)
	}
	r.precond, _ = NewMatchable(strings.TrimSpace(split[0]))
	r.postcond, _ = NewMatchable(strings.TrimSpace(split[1]))
	return nil
}

// [+stop+consonant+alveolar] > r / [+vowel+stress] _ [+vowel-stress]
// t d > r / [+vowel+stress] _ [+vowel-stress]
