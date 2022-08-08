package parser

import (
	"fmt"
	"strings"
)

type Rule interface {
    Apply(to string) string
}

type SimpleRule struct {
	from, to, env string
    precond, postcond string
}

func NewRule(rule string) Rule {
	r := &SimpleRule{}

	r.split(rule)

	return r
}

func (r *SimpleRule) split(rule string) {
	pointer := 0
	for i, char := range rule {
		switch char {
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

/*func (r *SimpleRule) parseEnv() {
    for i := 0; i < len(r.env); i++ {
        switch r.env[i] {
            case '_':
            
        }
    }
}*/

func (r *SimpleRule) Apply(to string) string {
    return replace(to, r.from, r.to)
}

/*
func trim(x string) string {
    lx := len(x)
    f, e := 0, lx - 1
    for ; f < lx && x[f] == ' '; f++ {}
    for ; e >= f && x[e] == ' '; e-- {}
    return x[f:e + 1]
}
*/

func replace(x string, find string, with string) string {
    ftil, cr, flen := 0, find[0], len(find)

    b := strings.Builder{}

    last := 0
    
    for i := 0; i < len(x); i++ {
        char := x[i]

        fmt.Println(char)
        if char == cr {
            ftil++
            if ftil >= flen {
                b.WriteString(x[last:i+1-flen])
                b.WriteString(with)
                ftil = 0
                last = i+1
            }
            cr = find[ftil]
        } else if ftil != 0 {
            ftil = 0
            cr = find[0]
        }
    }
    b.WriteString(x[last:])
    
    return b.String()
}

// [+stop+consonant+alveolar] > r / [+vowel+stress] _ [+vowel-stress]
// t d > r / [+vowel+stress] _ [+vowel-stress]