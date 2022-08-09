package parser

import (
	"fmt"
	"strings"
)

type Rule struct {
	from, to          Tree
	env               string
	precond, postcond Tree
}

func NewRule(rule string) *Rule {
	r := &Rule{}

	r.split(rule)

	r.parseEnv()

	return r
}

func (r *Rule) split(rule string) {
	pointer := 0
	var sto string
	for i := 0; i < len(rule); i++ {
		switch rule[i] {
		case '>':
			if r.from == nil {
				r.from, _ = setify(strings.TrimSpace(rule[pointer:i]))
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

	r.to, _ = setify(strings.TrimSpace(sto))
}

func (r *Rule) parseEnv() error {
	split := strings.SplitN(r.env, "_", 2)
	if len(split) != 2 {
		return fmt.Errorf("enviornment \"%s\" not in format \"{precondition} _ {postcondition}\"", r.env)
	}
	r.precond, _ = setify(strings.TrimSpace(split[0]))
	r.postcond, _ = setify(strings.TrimSpace(split[1]))
	return nil
}

func (r *Rule) Apply(text string) string {
	b := strings.Builder{}

	stages := [3]Tree{r.precond, r.from, r.postcond}
	cstage := 0

	var path []int
	p0, p1 := 0, 0
	lastWritten := 0

	// while in bounds of text
	for i := 0; i < len(text); {

	stageInc:
		if elmPath, mlen := MatchStart(stages[cstage], text[i:]); mlen != -1 {
			cstage++
			//fmt.Println(text[i:], i, cstage)

			i += mlen

			switch cstage {
			case 1:
				p0 = i
			case 2:
				p1 = i
				path = elmPath
			case 3:
				cstage = 0
				b.WriteString(text[lastWritten:p0])
				b.WriteString(string(FollowPath(r.to, path)))
				i = p1
				lastWritten = p1
			}

			if stages[cstage] == nil {
				goto stageInc
			}

		} else {
			cstage = 0
			i++
		}
	}

	b.WriteString(text[lastWritten:])

	return b.String()
}

/*
func (r *Rule) Apply(text string) string {
	b := strings.Builder{}

	fi := 0 // find index
	findArr := []Set{r.precond, r.from, r.postcond}
	stage := 0
	p0, p1 := 0, 0

	lastWritten := 0

	for i := 0; i < len(text); i++ {
		if len(findArr[stage]) == 0 || text[i] == findArr[stage][fi] {
			fi++
			if fi >= len(findArr[stage]) {
				// found!
			stageInc:
				stage++
				switch stage {
				case 1: // just found precond
					p0 = i + 1
				case 2: // just found "find"
					p1 = i + 1
				case 3:
					stage = 0
					b.WriteString(text[lastWritten:p0])
					b.WriteString(r.to)
					i = p1 - 1 // i = p1; i++ later, so -1
					lastWritten = p1
				}

				if len(findArr[stage]) == 0 {
					goto stageInc
				}

				fi = 0
			}
		}
	}

	b.WriteString(text[lastWritten:])
	return b.String()
}
*/

// [+stop+consonant+alveolar] > r / [+vowel+stress] _ [+vowel-stress]
// t d > r / [+vowel+stress] _ [+vowel-stress]
