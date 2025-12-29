# avro-phonetic-go

**Avro-style Banglish → বাংলা transliteration for modern Go applications**

<p align="center">
  <img src="avro-phonetic-go.png" alt="avro-phonetic-go logo" width="420">
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/mhshajib/avro-phonetic-go"><img src="https://pkg.go.dev/badge/github.com/mhshajib/avro-phonetic-go.png" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/mhshajib/avro-phonetic-go"><img src="https://goreportcard.com/badge/github.com/mhshajib/avro-phonetic-go?branch=main" alt="Go Report Card"></a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License: MIT">
  </a>
</p>

A Go library for converting Banglish (romanized Bangla) into Bangla script using an Avro Phonetic style grammar.

This project is inspired by the Avro Phonetic concept and implementation ideas demonstrated in the PHP library:

- Avro Phonetic: original keyboard concept and grammar
- avro-php by imerfanahmed (PHP implementation that influenced this design)

## Status

This repository ships with a compact, minimal default grammar for demonstration and correctness tests.

For production-grade transliteration parity, load a full grammar JSON using
`FromGrammarFile` or `FromGrammarReader`.

## Features

- Trie-based longest match scanning for fast conversion
- Context-aware rules (prefix and suffix constraints)
- Strict mode (Avro-compatible baseline)
- BD mode (opt-in shortcuts layered on top of strict behavior)
- Custom grammar support using JSON

## Installation

```bash
go get github.com/yourname/avro-phonetic-go
```

Update the module path to your GitHub org or user when publishing.

## Usage

### Strict mode (Avro-compatible)

```go
package main

import (
  "fmt"
  avrophonetic "github.com/yourname/avro-phonetic-go"
)

func main() {
  fmt.Println(avrophonetic.To("ami bangla gan gai"))
  // Output: আমি বাংলা গান গাই
}
```

Example file: `examples/strict/main.go`

### BD mode (opt-in shortcuts)

BD mode enables common Bangladeshi typing shortcuts while keeping strict grammar rules as the base.

```go
package main

import (
  "fmt"
  avrophonetic "github.com/yourname/avro-phonetic-go"
)

func main() {
  fmt.Println(avrophonetic.ToBD("tmi valo"))
  // Output: তুমি ভালো
}
```

Example file: `examples/bd/main.go`

## Using a custom grammar

### Load a grammar in strict mode

```go
g, err := avrophonetic.FromGrammarFile("./grammar.json")
if err != nil {
  panic(err)
}

a := avrophonetic.New(
  avrophonetic.WithGrammar(g),
  avrophonetic.Strict(),
)

fmt.Println(a.Parse("ami"))
```

### Load a grammar with BD mode enabled

BD mode applies BD-specific shortcuts before strict patterns.

```go
g, err := avrophonetic.FromGrammarFile("./grammar.json")
if err != nil {
  panic(err)
}

a := avrophonetic.New(
  avrophonetic.WithGrammar(g),
  avrophonetic.BDMode(),
)

fmt.Println(a.Parse("tmi valo"))
```

## Grammar JSON specification

### Base structure

```json
{
  "vowel": "aeiouAEIOU",
  "consonant": "bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ",
  "number": "0123456789",
  "casesensitive": "OI",
  "patterns": []
}
```

### Pattern object

```json
{
  "find": "ami",
  "replace": "আমি",
  "rules": []
}
```

Rules are optional.

## Rule system

Rules apply constraints based on surrounding characters.

### Rule fields

- `scope`  
  One of:

  - `vowel`, `!vowel`
  - `consonant`, `!consonant`
  - `punctuation`, `!punctuation`
  - `exact`, `!exact`

- `type`

  - `prefix` → checks character before the match
  - `suffix` → checks character after the match

- `value`  
  Required only for `exact` / `!exact`

---

## Grammar examples for all scopes

### 1. vowel (prefix)

Apply only if previous character is a vowel.

```json
{
  "find": "i",
  "replace": "ি",
  "rules": [{ "scope": "consonant", "type": "prefix" }]
}
```

### 2. !vowel (suffix)

Apply only if next character is not a vowel.

```json
{
  "find": "o",
  "replace": "ও",
  "rules": [{ "scope": "!vowel", "type": "suffix" }]
}
```

### 3. consonant (prefix)

```json
{
  "find": "a",
  "replace": "া",
  "rules": [{ "scope": "consonant", "type": "prefix" }]
}
```

### 4. !consonant (prefix)

```json
{
  "find": "n",
  "replace": "ন",
  "rules": [{ "scope": "!consonant", "type": "prefix" }]
}
```

### 5. punctuation (prefix)

Used to detect word boundaries.

```json
{
  "find": "ta",
  "replace": "টা",
  "rules": [{ "scope": "punctuation", "type": "prefix" }]
}
```

### 6. !punctuation (suffix)

```json
{
  "find": "na",
  "replace": "না",
  "rules": [{ "scope": "!punctuation", "type": "suffix" }]
}
```

### 7. exact (prefix)

Apply only if specific characters appear before the match.

```json
{
  "find": "ri",
  "replace": "ৃ",
  "rules": [{ "scope": "exact", "type": "prefix", "value": "k" }]
}
```

### 8. !exact (suffix)

Apply only if the next characters are not a specific value.

```json
{
  "find": "e",
  "replace": "ে",
  "rules": [{ "scope": "!exact", "type": "suffix", "value": "r" }]
}
```

---

## Compatibility notes

- The bundled grammar is intentionally minimal.
- Full Avro parity depends on the grammar JSON you provide.
- The rule engine is designed to support typical Avro-style constraints.
- Edge-case equivalence depends on grammar completeness.

## Roadmap

### v0.1.x — Grammar-driven engine (current)

- Trie-based longest-match parser
- Context-aware rule engine (prefix/suffix)
- Strict (Avro-style) mode
- BD mode (opt-in shortcuts)
- Custom grammar loading via JSON

### v0.2.x — Full Avro grammar parity (planned)

The goal of this milestone is to reach practical parity with the commonly used
Avro Phonetic grammar.

Planned work includes:

- Porting the complete Avro grammar into JSON form
- Adding comprehensive vowel-kar and conjunct handling
- Expanding rule coverage for edge cases
- Introducing golden test cases to validate output compatibility
- Documenting known behavioral differences, if any

Parity will be defined by **output equivalence given the same grammar input**,
not by internal implementation details.

Contributions in the form of grammar rules, test cases, and validation examples
are welcome.

## Benchmarks

```bash
go test ./... -bench=.
```

## Credits

- Avro Phonetic keyboard: original idea and grammar concept
- PHP implementation reference: `https://github.com/imerfanahmed/avro-php`

This project is an independent Go implementation and is not affiliated with the original Avro project.

## Contributing

Issues and pull requests are welcome at:

https://github.com/mhshajib/avro-phonetic-go

When contributing code, please ensure that you follow standard Go formatting and conventions.

Please include tests for any functional changes.

## Contributors

<a href="https://github.com/mhshajib/avro-phonetic-go/graphs/contributors">
  <img src="https://img.shields.io/github/contributors/mhshajib/avro-phonetic-go?style=for-the-badge" alt="Contributors"/>
</a>

[![Contributors](https://contrib.rocks/image?repo=mhshajib/avro-phonetic-go&max=36&v=2)](https://github.com/mhshajib/avro-phonetic-go/graphs/contributors)

## License

MIT. See `LICENSE`.
