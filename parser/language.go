package parser

import "fmt"

type Language struct{}

func NewLanguage(rule, txt string) (l *Language) {

	r1 := NewRule(rule)
	fmt.Println(r1.Apply(txt))

	return &Language{}
}

/*

reserved chars:
> / _ { } [ ] + - ( ) #
    maybe ? * = !

implemented:
> / _ { } ( ) #

not implemented:
[ ] + -
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
[ ] [a] > b
[ ] [a] > [b]
    // provided len(a) == len(b)

Arithmitic with named sets intrasectionally
[ ] [a+b] > c
[ ] [a+b+c] > d
[ ] [a-b] > c
[ ] [a+b-c] > d
[ ] [a-b-c] > d
[ ] [a-(b-c)] > d

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
