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

	features := map[string]*ValueSet{}
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
	vs := &ValueSet{}
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


@manner stop = p b t d k g
@manner nasal = m n ng
voice = b m d g n ng
@place bilabial = m p b
@place alveolar = n t d

[STOP = +stop-labial]{a {i u} {e o}} > [STOP -> +labial]{æ y œ} / {e i} _ [+stop]

non-labial stops and a vowel to a labial and "a" when preceded by "e" and succeeded by a stop

Sequence{
    Set{label: STOP, set: difference(union(None, stop), labial)},
    Set{label: None, set:[a e i o u]}
}

TO

Sequence{
    Term{from: STOP, set: }
}
*/
