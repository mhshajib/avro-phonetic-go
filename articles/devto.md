---
title: "avro-phonetic-go: Avro-style Banglish to বাংলা transliteration in Go"
published: false
description: "A trie-based, context-aware Avro Phonetic transliteration engine in Go with strict compatibility and an opt-in BD mode."
tags: go, bangla, transliteration, opensource
---

If you work with Bengali users, you already know the problem.

People type Bangla using Latin letters. Not carefully. Not consistently.
Yet they expect your application to understand them.

The Avro Phonetic keyboard showed that this problem can be solved using
a grammar-driven approach: pattern matching combined with local context rules.

This article introduces **avro-phonetic-go**, a Go library that implements
an Avro-style phonetic transliteration engine with a clean, production-oriented API.

## Design goals

The library was built with the following principles:

- Grammar-driven, not hardcoded
- Deterministic output using longest-match scanning
- Explicit rule evaluation (prefix and suffix constraints)
- No magic heuristics in strict mode
- Optional BD mode for modern Bangladeshi typing shortcuts

The overall idea is credited to the Avro Phonetic keyboard.
The internal design of this library was strongly influenced by
the PHP reference implementation `imerfanahmed/avro-php`.

## Quick example

Strict mode:

```go
fmt.Println(avrophonetic.To("ami bangla gan gai"))
// আমি বাংলা গান গাই
```

BD mode (opt-in):

```go
fmt.Println(avrophonetic.ToBD("tmi valo"))
// তুমি ভালো
```

## How it works

At a high level, the engine works in four steps:

1. Load a grammar consisting of patterns and rules
2. Build a trie for fast longest-match lookup
3. Scan the input left to right
4. Validate candidate patterns using local context rules

Rules are evaluated only against immediate neighbors.
There is no backtracking, which keeps the algorithm predictable and fast.

## Strict mode vs BD mode

Strict mode is intended to behave as a clean Avro-style baseline.
If you want exact behavior, use strict mode only.

BD mode layers a small set of additional patterns on top of strict grammar.
These patterns capture real-world Bangladeshi typing habits
without polluting the base grammar.

This separation keeps the engine usable for both compatibility
and user-experience–oriented applications.

## Custom grammar support

The engine is fully grammar-driven.

You can load a complete grammar JSON file to reach parity with the grammar you prefer:

```go
g, _ := avrophonetic.FromGrammarFile("grammar.json")
a := avrophonetic.New(avrophonetic.WithGrammar(g))
fmt.Println(a.Parse("ami"))
```

This makes the library suitable for search indexing, chat processing,
form normalization, and NLP pipelines.

## Credits

- Avro Phonetic keyboard: original idea and grammar concept
- PHP reference implementation: `imerfanahmed/avro-php`

This project is an independent Go implementation and is not affiliated
with the original Avro project.

## Closing thoughts

Transliteration is a small feature with outsized impact.

If your product handles Bengali user input,
doing this well is not optional anymore.

Repository: avro-phonetic-go
