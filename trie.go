package avrophonetic

type trieNode struct {
	next     map[rune]*trieNode
	patterns []Pattern // patterns that share the same Find string
	terminal bool
}

type trie struct {
	root *trieNode
}

func newTrie() *trie {
	return &trie{root: &trieNode{next: make(map[rune]*trieNode)}}
}

func (t *trie) insert(find string, p Pattern) {
	n := t.root
	for _, r := range []rune(find) {
		if n.next == nil {
			n.next = make(map[rune]*trieNode)
		}
		if _, ok := n.next[r]; !ok {
			n.next[r] = &trieNode{next: make(map[rune]*trieNode)}
		}
		n = n.next[r]
	}
	n.terminal = true
	n.patterns = append(n.patterns, p)
}

// matchLongest finds the longest match starting at runes[pos].
func (t *trie) matchLongest(runes []rune, pos int) (matchedLen int, patterns []Pattern) {
	n := t.root
	bestLen := 0
	var best []Pattern

	for i := pos; i < len(runes); i++ {
		r := runes[i]
		next, ok := n.next[r]
		if !ok {
			break
		}
		n = next
		if n.terminal {
			bestLen = i - pos + 1
			best = n.patterns
		}
	}
	return bestLen, best
}
