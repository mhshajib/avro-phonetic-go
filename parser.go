package avrophonetic

import "unicode"

// Parser converts input using a compiled trie and rule evaluation.
type Parser struct {
	g    Grammar
	trie *trie
	set  runeSets
}

func NewParser(g Grammar) *Parser {
	t := newTrie()
	for _, p := range g.Patterns {
		if p.Find == "" {
			continue
		}
		t.insert(p.Find, p)
	}
	return &Parser{
		g:    g,
		trie: t,
		set:  newRuneSets(g),
	}
}

func (p *Parser) Parse(text string) string {
	if text == "" {
		return ""
	}
	in := []rune(text)
	out := make([]rune, 0, len(in))

	i := 0
	for i < len(in) {
		ml, patterns := p.trie.matchLongest(in, i)
		if ml == 0 {
			out = append(out, in[i])
			i++
			continue
		}
		// Choose the first pattern whose rules match; otherwise fallback to the first pattern.
		chosen := patterns[0]
		for _, cand := range patterns {
			if p.rulesMatch(in, i, ml, cand.Rules) {
				chosen = cand
				break
			}
		}
		out = append(out, []rune(chosen.Replace)...)
		i += ml
	}
	return string(out)
}

func (p *Parser) rulesMatch(in []rune, pos, matchedLen int, rules []Rule) bool {
	if len(rules) == 0 {
		return true
	}
	for _, r := range rules {
		if !p.ruleMatch(in, pos, matchedLen, r) {
			return false
		}
	}
	return true
}

func (p *Parser) ruleMatch(in []rune, pos, matchedLen int, r Rule) bool {
	kind := r.Scope
	neg := false
	if len(kind) > 0 && kind[0] == '!' {
		neg = true
		kind = kind[1:]
	}
	ok := false

	switch kind {
	case "vowel":
		ok = p.checkNeighborSet(in, pos, matchedLen, r.Type, p.set.isVowel)
	case "consonant":
		ok = p.checkNeighborSet(in, pos, matchedLen, r.Type, p.set.isConsonant)
	case "punctuation":
		ok = p.checkNeighborSet(in, pos, matchedLen, r.Type, isPunctOrSpace)
	case "exact":
		ok = p.checkNeighborExact(in, pos, matchedLen, r.Type, []rune(r.Value))
	default:
		// Unknown rule types are treated as non-matching to keep behavior explicit.
		ok = false
	}

	if neg {
		return !ok
	}
	return ok
}

func (p *Parser) checkNeighborSet(in []rune, pos, matchedLen int, ruleType string, fn func(rune) bool) bool {
	switch ruleType {
	case "prefix":
		ix := pos - 1
		if ix < 0 {
			return false
		}
		return fn(in[ix])
	case "suffix":
		ix := pos + matchedLen
		if ix >= len(in) {
			return false
		}
		return fn(in[ix])
	default:
		return false
	}
}

func (p *Parser) checkNeighborExact(in []rune, pos, matchedLen int, ruleType string, val []rune) bool {
	if len(val) == 0 {
		return false
	}
	switch ruleType {
	case "prefix":
		start := pos - len(val)
		if start < 0 {
			return false
		}
		for i := range val {
			if in[start+i] != val[i] {
				return false
			}
		}
		return true
	case "suffix":
		start := pos + matchedLen
		end := start + len(val)
		if start < 0 || end > len(in) {
			return false
		}
		for i := range val {
			if in[start+i] != val[i] {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func isPunctOrSpace(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}
	// Treat common punctuation as boundaries. Extendable by users.
	return unicode.IsPunct(r) || r == '।' || r == '—' || r == '–'
}
