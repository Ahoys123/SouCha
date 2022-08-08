package parser

import "fmt"

type RuleContext struct {
    sth map[string]GroupHash
    hts map[GroupHash]string

    labels map[string]uint8
}

func Hash(x []bool) GroupHash {
    var n GroupHash = 0
    for i, v := range x {
        if v {
            n += 1 << i
        }
    }
    return n
}

func SetFlag(n GroupHash, pos uint8, to bool) GroupHash {
    if to {
        n |= (1 << pos)
    } else {
        n &= ^(1 << pos)
    }
    return n
}

func GetFlag(n GroupHash, pos uint8, to bool) GroupHash {
    return (n & (1 << pos)) >> pos
}

func NewRuleContext(groups map[string][]string) (*RuleContext, error) {
    sth := make(map[string]GroupHash, 0)
    hts := make(map[GroupHash]string, 0)

    labels := make(map[string]uint8, len(groups))
    var n uint8 = 0
    for name, vals := range groups {
        labels[name] = n
        for _, p := range vals {
            mask, ok := sth[p]
            if !ok {
                mask = 0
            }
            sth[p] = SetFlag(mask, n, true)
        }
        n++
    }

    var err error
    for k, v := range sth {
        p, notUnique := hts[v]
        if notUnique && err == nil {
            err = fmt.Errorf("phoneme %s has the same mask as %s; 0b%04b\n", p, k, v)
        }
        hts[v] = k
    }

    return &RuleContext{sth, hts, labels}, err
}

type GroupHash int

/*

Base
[X] a > b 
[X] a >  
    // a gets deleted everywhere

Enviornments
[ ] a > b / c _ d

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