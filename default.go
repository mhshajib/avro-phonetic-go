package avrophonetic

import "sync"

var (
	defaultOnce sync.Once
	defaultAvro *Avro

	defaultBDOnce sync.Once
	defaultBDAvro *Avro
)

// Default returns a singleton converter in strict mode.
func Default() *Avro {
	defaultOnce.Do(func() {
		defaultAvro = New(Strict())
	})
	return defaultAvro
}

// DefaultBD returns a singleton converter in BD mode.
func DefaultBD() *Avro {
	defaultBDOnce.Do(func() {
		defaultBDAvro = New(BDMode())
	})
	return defaultBDAvro
}
