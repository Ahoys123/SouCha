package parser

import "testing"

func TestSetArithmetic(t *testing.T) {
	l := NewLanguage(`stop = p t k b d g
	consonant = p t k b d g m n j
	vowel = a e i o u
	alveolar = t d n
	[-alveolar+stop] > p
	[+alveolar+stop] > k`)
	if done := l.Evolve("nate"); done != "nake" {
		t.Errorf("rule [+alveolar+stop] > k not obeyed! Got nate > %s", done)
	}

	if done := l.Evolve("nake"); done != "nape" {
		t.Errorf("rule [-alveolar+stop] > p not obeyed! Got nake > %s", done)
	}

	l = NewLanguage(`stop = p t k b d g
	consonant = p t k b d g m n j
	vowel = a e i o u
	alveolar = t d n
	[+alveolar|stop] > p`)

	if done := l.Evolve("nane"); done != "pape" {
		t.Errorf("rule [+alveolar|stop] > p not obeyed! Got nane > %s", done)
	}

	if done := l.Evolve("tame"); done != "pame" {
		t.Errorf("rule [+alveolar|stop] > p not obeyed! Got tame > %s", done)
	}

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
