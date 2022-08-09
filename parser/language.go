package parser

import "fmt"

type Language struct{}

func NewLanguage(test string) (l *Language) {
	r := NewRule("日 > 火 / 喪 _ 喪")
	fmt.Println(r.Apply("喪日喪日喪"))

	x, _ := setify("{{aaa bbb} ccc}")
	fmt.Println(x)
	return &Language{}
}

/*

reserved chars:
> / _ { } [ ] + - ( )
    maybe ? * =

Base
[X] a > b
[X] a >
    // a gets deleted everywhere

Enviornments
[X] a > b / c _ d

Basic unnamed sets
[ ] a b > c
[ ] a b > c d
[ ] {a b} c > d e
[ ] {a {b c} d} e > {f g h} i

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
