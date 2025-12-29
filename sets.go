package avrophonetic

type runeSets struct {
	vowelSet     map[rune]struct{}
	consonantSet map[rune]struct{}
}

func newRuneSets(g Grammar) runeSets {
	rs := runeSets{
		vowelSet:     make(map[rune]struct{}),
		consonantSet: make(map[rune]struct{}),
	}
	for _, r := range []rune(g.Vowel) {
		rs.vowelSet[r] = struct{}{}
	}
	for _, r := range []rune(g.Consonant) {
		rs.consonantSet[r] = struct{}{}
	}
	return rs
}

func (s runeSets) isVowel(r rune) bool {
	_, ok := s.vowelSet[r]
	return ok
}

func (s runeSets) isConsonant(r rune) bool {
	_, ok := s.consonantSet[r]
	return ok
}
