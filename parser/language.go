package parser

import "fmt"

type Language struct {
    ctx *RuleContext
}

func NewLanguage(test string) (l *Language) {
    r := NewRule("xa >  / c")
    fmt.Println(r.Apply("baxmaxa"))
    /*
    ctx, err := NewRuleContext(map[string][]string{
        //"bilabial" : {"p", "b"},
        "stop" : {"p", "t", "k"},
        //"alveolar" : {"t"},
    })

    // p > [-stop]
    // makes no sense if p has multiple non stopped variants (m & f, for example)'
    // p > [+alveolar] also means taking out [+bilabial] flag from p
    // 11010 -> 11011 -> 11001
    // 11010 -> 11011 -> 10011? X because not possible to have bilabial & alveolar flags set at same time
    
    if err != nil {
        fmt.Println(err)
    }

    for k, v := range ctx.hts {
        fmt.Printf("hash: %04b, string: %s\n", k, v)
    }
    return l*/
    return &Language{}
}