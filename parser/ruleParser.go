package parser

import "fmt"

type Rule interface {
    Apply(to []rune) []rune
}

type SimpleRule struct {
	from, to, env []rune
}

func NewRule(rule []rune) Rule {
	r := &SimpleRule{}

	r.split(rule)

	return r
}

func (r *SimpleRule) split(rule []rune) {
	pointer := 0
	for i, char := range rule {
		switch char {
		case '>':
			if r.from == nil {
				r.from = trim(rule[pointer:i])
			}
			pointer = i + 1
		case '/':
			r.to = rule[pointer:i]
			pointer = i + 1
		}
	}

	if r.to == nil {
		r.to = rule[pointer:]
	} else {
		r.env = trim(rule[pointer:])
	}

    r.to = trim(r.to)
}

func (r *SimpleRule) Apply(to []rune) []rune {
    return replace(to, r.from, r.to)
}

func trim(x []rune) []rune {
    lx := len(x)
    f, e := 0, lx - 1
    for ; f < lx && x[f] == ' '; f++ {}
    for ; e >= f && x[e] == ' '; e-- {}
    return x[f:e + 1]
}

func replace(x []rune, find []rune, with []rune) []rune {
    ftil, cr, flen, wlen := 0, find[0], len(find), len(with)
    for i := 0; i <= len(x) - flen; i++ {
        char := x[i]
        
        fmt.Println("\n", x, i)
        if char == cr {
            ftil++
            if ftil >= flen {
                x = indexReplace(x, i, flen, with, wlen)
                i += wlen - flen
                ftil = 0
            }
            cr = find[ftil]
        }
    }
    return x
}

func indexReplace(x []rune, i int, flen int, with []rune, wlen int) []rune {
    // if 1 : 1 correspondence, replace letters directly
    if wlen == flen {
        for j := 0; j < flen; j++ {
            x[i + j] = with[j]
        }
        return x
    }

    // if not, must allocate another slice
    fmt.Println(x[:i], with, x[i+flen:])
    
    b := append(append(x[:i], with...), x[i+flen:]...) 
    fmt.Println(string(b))
    return b
    // x[:i] + with + x[i+wlen:]
}

// [+stop+consonant+alveolar] > r / [+vowel+stress] _ [+vowel-stress]
// t d > r / [+vowel+stress] _ [+vowel-stress]