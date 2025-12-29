package avrophonetic

import "testing"

func TestStrictExamples(t *testing.T) {
	got := To("ami bangla gan gai")
	want := "আমি বাংলা গান গাই"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestBDExamples(t *testing.T) {
	got := ToBD("tmi valo")
	want := "তুমি ভালো"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCustomGrammar(t *testing.T) {
	g := Grammar{
		Vowel:     "aeiou",
		Consonant: "bcdfghjklmnpqrstvwxyz",
		Patterns: []Pattern{
			{Find: "hello", Replace: "হ্যালো"},
		},
	}
	a := New(WithGrammar(g))
	got := a.Parse("hello")
	if got != "হ্যালো" {
		t.Fatalf("got %q", got)
	}
}
