package fnttw

import (
	"reflect"
	"testing"
)

func Range(n int) []uint64 {
	a := make([]uint64, n)
	for i := range a {
		a[i] = uint64(i)
	}
	return a
}

func TestNtt(t *testing.T) {
	N := 1 << 10
	a := Range(N)
	// TODO(nsamar): Clean up arguments, they are messy
	q := NttFriendlyPrimes(N, 1, IntLog2(N)+4)[0]
	b := make([]uint64, len(a))
	modulus := MakeMontgomery64BitModulus(q)
	// modulus := Simple64BitModulus{q}
	copy(b, a)
	Ntt(a, modulus)
	TrivialNtt(b, modulus.Modulus())
	if !reflect.DeepEqual(a, b) {
		t.Error("NTT failed\n", a, "\n", b)
	}
	if InverseTrivialNtt(b, modulus.Modulus()); !reflect.DeepEqual(b, Range(N)) {
		t.Error("Trivial InverseNTT failed\n", b, "\n", Range(N))
	}
	if InverseNtt(a, modulus); !reflect.DeepEqual(a, Range(N)) {
		t.Error("Inverse NTT failed", a, Range(N))
	}
}
