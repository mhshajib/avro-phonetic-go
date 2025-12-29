// Package avrophonetic provides an Avro-style phonetic transliteration engine
// for converting Banglish (romanized Bangla) text into Bangla Unicode script.
//
// The engine is grammar-driven and uses a trie-based longest-match scanner
// combined with context-aware rules (prefix and suffix constraints).
//
// Two operating modes are supported:
//
//   - Strict mode: Avro-compatible baseline behavior
//   - BD mode: Optional Bangladeshi typing shortcuts layered on top of strict mode
//
// The default embedded grammar is intentionally minimal. For production usage
// and full transliteration parity, users are expected to load a complete
// Avro-compatible grammar JSON using FromGrammarFile or FromGrammarReader.
//
// This package is an independent Go implementation inspired by the Avro
// Phonetic concept and the PHP reference implementation by imerfanahmed.
package avrophonetic

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// To converts romanized Bangla (Banglish) input into Bangla using the embedded default grammar
// in strict (Avro-compatible) mode.
func To(text string) string {
	return Default().Parse(text)
}

// ToBD converts input using the embedded default grammar and BD-tuned mode enabled.
func ToBD(text string) string {
	return DefaultBD().Parse(text)
}

type Mode int

const (
	ModeStrict Mode = iota
	ModeBD
)

// Avro is the main converter object. It holds a compiled trie and a grammar.
type Avro struct {
	parser *Parser
}

// New returns a converter using the embedded grammar (strict mode by default).
func New(opts ...Option) *Avro {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	g := cfg.grammar
	if cfg.mode == ModeBD {
		g = MergeGrammar(BDExtras(), g)
	}
	return &Avro{parser: NewParser(g)}
}

// Parse converts the input string to Bangla.
func (a *Avro) Parse(text string) string {
	if a == nil || a.parser == nil {
		return text
	}
	return a.parser.Parse(text)
}

// Parser exposes the underlying parser (advanced usage).
func (a *Avro) Parser() *Parser {
	if a == nil {
		return nil
	}
	return a.parser
}

type config struct {
	mode    Mode
	grammar Grammar
}

func defaultConfig() config {
	return config{
		mode:    ModeStrict,
		grammar: DefaultGrammar(),
	}
}

type Option func(*config)

// Strict forces strict mode output (Avro-compatible behavior).
func Strict() Option {
	return func(c *config) { c.mode = ModeStrict }
}

// BDMode enables BD-tuned shortcuts on top of strict grammar.
func BDMode() Option {
	return func(c *config) { c.mode = ModeBD }
}

// WithGrammar uses a caller-provided grammar instead of the embedded default.
func WithGrammar(g Grammar) Option {
	return func(c *config) { c.grammar = g }
}

// FromGrammarReader loads grammar JSON from an io.Reader.
func FromGrammarReader(r io.Reader) (Grammar, error) {
	if r == nil {
		return Grammar{}, errors.New("nil reader")
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return Grammar{}, err
	}
	var g Grammar
	if err := json.Unmarshal(b, &g); err != nil {
		return Grammar{}, err
	}
	return g, nil
}

// FromGrammarFile loads grammar JSON from a file path.
func FromGrammarFile(path string) (Grammar, error) {
	f, err := os.Open(path)
	if err != nil {
		return Grammar{}, err
	}
	defer f.Close()
	return FromGrammarReader(f)
}
