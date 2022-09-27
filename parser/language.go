package parser

import (
	"fmt"
	"strings"
)

type Language struct {
	rules []*Rule
	ctx   *RuleContext
}

func NewLanguage(txt string) (l *Language) {

	l = &Language{}

	features := map[string]ValueSet{}
	rules := []string{}

	// Create feature sets from context & store rules as strings
	for i, line := range strings.Split(txt, "\n") {
		var v []string
		if v = strings.SplitN(line, "=", 2); len(v) > 1 {
			vs, _ := NewValueSet(v[1])
			features[strings.TrimSpace(v[0])] = vs
		} else if strings.ContainsRune(line, '>') {
			rules = append(rules, line)
		} else if !strings.Contains(line, "//") && !(len(strings.TrimSpace(line)) == 0) {
			fmt.Printf("couldn't parse line %d, \"%s\"\n", i, line)
		}
	}

	// Create Universal set from union of all other sets
	vs := ValueSet{}
	for _, v := range features {
		vs = vs.Union(v)
	}

	l.ctx = &RuleContext{
		features,
		vs,
	}

	// convert string rules to *Rule types, AFTER context has been established
	for _, v := range rules {
		l.rules = append(l.rules, l.NewRule(v))
	}

	return l
}

func (l *Language) Evolve(txt string) string {
	txt = " " + txt + " "
	for _, rule := range l.rules {
		txt = rule.Apply(txt)

		fmt.Println(txt)
	}
	return strings.TrimSpace(txt)
}

/*

reserved chars:
> / _ { } [ ] + - ( ) # ! | &

implemented:
> / _ { } ( ) # [ ] + - ! | &

not implemented:
    + !
    ! in env / outside []
    + for 1 or more of

TODO:
    multiple envs
        allow for / (for multiple envs) and ! (for not envs)
    more waiters
        allow for + to wait on last char
    backtracking on precondition sets
        if one precond matches but from/postcond fails, try finding another in precond before continuing
    assimilation
        allow a "to" arg to take a certain feature from whatever "from" matches
            could be [+nasal] > [@place] / _ [+stop@place]
            @place = {bilabial alveolar velar}
    simpler [!{phone}] notation?

Base
[X] a > b
[X] a >
    // a gets deleted everywhere

Enviornments
[X] a > b / c _ d

Basic unnamed sets
[X] a b > c
[X] a b > c d
[X] {a b} c > d e
[X] {a {b c} d} e > {f g h} i

Basic named sets
[X] [a] > b
[?] [a] > [b]
    // provided len(a) == len(b)

Arithmitic with named sets intrasectionally
[X] [a+b] > c
[X] [a+b+c] > d
[X] [a-b] > c
[X] [a+b-c] > d
[X] [a-b-c] > d
[X] [a-(b-c)] > d

Arithmitic with named sets intersectionally (see Implementation)
[ ] [a] > *[-a]
*/

/*

Implementation of intersectional named sets

stops = p b t d k g
voice = m b n d ng g
labial = p b m
alveolar = t d n
velar = k g ng
    stops voice labial alveolar velar
p = 1     0     1      0        0
b = 1     1     1      0        0
m = 0     1     1      0        0


p t k > *[+voice]
if p:
    newMask = my mask
    change [voice] flag on mask
    look at all other masks, find closest related one

*/
