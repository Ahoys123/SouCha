# SouCha

SouCha is a Go package that allows a user to simulate historical sound changes. Provided a list of sound changes ("rules"), SouCha can apply each change sequentally to a set of words or phrases. 

## Usage

### In Code

All logic is contained in package parser. Create a new language with `l := parser.NewLanguage("__rules to parse__")`. Then you can evolve a segment with `l.Evolve("word/phrase")`.

### In Language File

#### Rules
Rules must be in the form `initialSound > finalSound / preCondition _ postCondition`. If conditions for the sound change don't exist, then you can omit everything after the slash. Any component can be omitted if it is the empty set ∅ (it is non-restrictive).

Any of these four components can contain:
* Spaces and commas, indicating a break between two parts of a set: `x y` and `x,y` both indicate the set of x or y.
   * To indicate a sequence of consecutive phonemes, put no spaces in between: `xy` indicates xy in sequence.
   * These consecutive phonemes can also be sets of phonemes: `x{y z}` indicates xy or xz.
* Curly braces `{}`, indicating any one of a set: `{x y}` allows you to group phonemes.
* Parentheses `()`, indicating an optional phoneme: `(ʰ)` indicates either ʰ or nothing.
* Square brackets `[]`, indicating a combination of named sets, with their own rules inside:
   * `+setname` and `&setname` indicate an intersection, meaning the matched phoneme must be in that set.
   * `-setname` and `!setname` indicate the complement, meaning the matched phoneme cannot be in that set.
   * `|setname` indicates a union, meaning the matched phoneme can be in this set or the set before it.
   * `()` indicates logical groupings, NOT optionality– use `|setname` for optionality instead.
   * `{}` indicates an anonymous set; it acts as a named set outside of the braces, but the symbols inside are assumed to be phonemes, not named sets.
   * `NAME:` at the beginning of square brackets indicates the set's name. It can be any token with no spaces and referenced later using `[NAME:]` like any other phoneme to describe the specific phoneme matched previously. 

`#` is a special character in `preCondition` and `postCondition`, indicating a word boundary (beginning/end of a word).

#### Named Sets
A named set must be in the form `setname = x y z`, either space or comma-seperated. They cannot contain other named sets, curly braces, or parentheses as of now. They represent a set of phonemes.

## Motivation
There are many sound change appliers online already, but SouCha allows for a rule input that is closer to general linguistic conventions. Rules might look like 

`[V:+vowel] > / [V:] _ {x y}`

I.e, vowels disappear when the same vowel appears before and is followed by /x/ or /y/.
