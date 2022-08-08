package parser

import (
	"fmt"
	"strings"
)

type Rule struct {
	from, to, env     string
	precond, postcond string
}

func NewRule(rule string) *Rule {
	r := &Rule{}

	r.split(rule)

	r.parseEnv()

	return r
}

func (r *Rule) split(rule string) {
	pointer := 0
	for i := 0; i < len(rule); i++ {
		switch rule[i] {
		case '>':
			if r.from == "" {
				r.from = strings.TrimSpace(rule[pointer:i])
			}
			pointer = i + 1
		case '/':
			r.to = rule[pointer:i]
			pointer = i + 1
		}
	}

	if r.to == "" {
		r.to = rule[pointer:]
	} else {
		r.env = strings.TrimSpace(rule[pointer:])
	}

	r.to = strings.TrimSpace(r.to)
}

func (r *Rule) parseEnv() error {
	split := strings.SplitN(r.env, "_", 2)
	if len(split) != 2 {
		return fmt.Errorf("enviornment \"%s\" not in format \"{precondition} _ {postcondition}\"", r.env)
	}
	r.precond, r.postcond = strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
	return nil
}

func (r *Rule) Apply(text string) string {
	b := strings.Builder{}

	fi := 0 // find index
	findArr := []string{r.precond, r.from, r.postcond}
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

// [+stop+consonant+alveolar] > r / [+vowel+stress] _ [+vowel-stress]
// t d > r / [+vowel+stress] _ [+vowel-stress]
