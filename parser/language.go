package parser

import (
	"strings"
)

type Language struct {
	rules []*Rule
	ctx   *RuleContext
}

func NewLanguage(context string, rules string) (l *Language) {

	l = &Language{}

	features := map[string]ValueSet{}

	for _, line := range strings.Split(context, "\n") {
		v := strings.SplitN(line, "=", 2)
		vs, _ := NewValueSet(v[1])
		features[strings.TrimSpace(v[0])] = vs
	}

	vs := ValueSet{}
	for _, v := range features {
		vs = vs.Union(v)
	}

	l.ctx = &RuleContext{
		features,
		vs,
	}

	for _, v := range strings.Split(rules, "\n") {
		l.rules = append(l.rules, l.NewRule(v))
	}

	return l
}

func (l *Language) Apply(txt string) string {
	for _, rule := range l.rules {
		txt = rule.Apply(txt)
	}
	return txt
}

/*

reserved chars:
> / _ { } [ ] + - ( ) #
    maybe ? * = !

implemented:
> / _ { } ( ) # [ ] + -

not implemented:
    ? * = !

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
