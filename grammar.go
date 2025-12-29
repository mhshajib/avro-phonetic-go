package avrophonetic

import "strings"

// Grammar defines the rule sets and patterns used by the parser.
type Grammar struct {
	Patterns      []Pattern `json:"patterns"`
	Vowel         string    `json:"vowel"`
	Consonant     string    `json:"consonant"`
	Number        string    `json:"number"`
	CaseSensitive string    `json:"casesensitive"`
}

type Pattern struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
	Rules   []Rule `json:"rules,omitempty"`
}

// Rule describes a context constraint.
type Rule struct {
	Scope string `json:"scope"` // vowel, consonant, punctuation, exact, or negations: !vowel, !consonant, !punctuation, !exact
	Type  string `json:"type"`  // prefix or suffix
	Value string `json:"value,omitempty"`
}

// DefaultGrammar returns a compact embedded grammar.
// This repo intentionally ships a minimal baseline grammar to keep the library small and auditable.
// You can load a full grammar via FromGrammarFile/FromGrammarReader.
func DefaultGrammar() Grammar {
	return Grammar{
		Vowel:         "aeiouAEIOU",
		Consonant:     "bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ",
		Number:        "0123456789",
		CaseSensitive: "OI",
		Patterns: []Pattern{
			// Core examples from Avro documentation and common Banglish.
			{Find: "ami", Replace: "আমি"},
			{Find: "bangla", Replace: "বাংলা"},
			{Find: "gan", Replace: "গান"},
			{Find: "gai", Replace: "গাই"},
			{Find: "tumi", Replace: "তুমি"},
			{Find: "tomra", Replace: "তোমরা"},
			{Find: "kisu", Replace: "কিছু"},
			{Find: "na", Replace: "না"},
			// A small set of letters; a full grammar can be swapped in.
			{Find: "a", Replace: "া", Rules: []Rule{{Scope: "consonant", Type: "prefix"}}}, // consonant + a => aa kar
			{Find: "i", Replace: "ি", Rules: []Rule{{Scope: "consonant", Type: "prefix"}}},
			{Find: "u", Replace: "ু", Rules: []Rule{{Scope: "consonant", Type: "prefix"}}},
			{Find: "e", Replace: "ে", Rules: []Rule{{Scope: "consonant", Type: "prefix"}}},
			{Find: "o", Replace: "ো", Rules: []Rule{{Scope: "consonant", Type: "prefix"}}},
			{Find: "a", Replace: "অ"}, // standalone
			{Find: "i", Replace: "ই"},
			{Find: "u", Replace: "উ"},
			{Find: "e", Replace: "এ"},
			{Find: "o", Replace: "ও"},
			{Find: "k", Replace: "ক"},
			{Find: "g", Replace: "গ"},
			{Find: "n", Replace: "ন"},
			{Find: "m", Replace: "ম"},
			{Find: "t", Replace: "ত"},
			{Find: "b", Replace: "ব"},
			{Find: "l", Replace: "ল"},
			{Find: "y", Replace: "য়"},
		},
	}
}

// BDExtras defines additional patterns for BD-tuned mode.
// These are applied before strict patterns (higher priority).
func BDExtras() Grammar {
	g := DefaultGrammar()
	g.Patterns = []Pattern{
		{Find: "tmi", Replace: "তুমি"},
		{Find: "tmra", Replace: "তোমরা"},
		{Find: "kmn", Replace: "কেমন"},
		{Find: "valo", Replace: "ভালো"},
		{Find: "bhalo", Replace: "ভালো"},
		{Find: "jodi", Replace: "যদি"},
		{Find: "kisuina", Replace: "কিছুই না"},
		{Find: "nai", Replace: "নাই"},
		{Find: "hoy", Replace: "হয়"},
		{Find: "hoyeche", Replace: "হয়েছে"},
		{Find: "ta", Replace: "টা", Rules: []Rule{{Scope: "punctuation", Type: "prefix"}}}, // word boundary "ta" => "টা"
	}
	// keep sets; patterns only matter
	return g
}

// MergeGrammar returns a merged grammar where "front" patterns take precedence over "back".
func MergeGrammar(front, back Grammar) Grammar {
	out := back
	out.Patterns = append([]Pattern{}, front.Patterns...)
	out.Patterns = append(out.Patterns, back.Patterns...)

	// Prefer explicit set values if front provided them.
	if strings.TrimSpace(front.Vowel) != "" {
		out.Vowel = front.Vowel
	}
	if strings.TrimSpace(front.Consonant) != "" {
		out.Consonant = front.Consonant
	}
	if strings.TrimSpace(front.Number) != "" {
		out.Number = front.Number
	}
	if strings.TrimSpace(front.CaseSensitive) != "" {
		out.CaseSensitive = front.CaseSensitive
	}
	return out
}
