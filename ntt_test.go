package fnttw

import (
	"reflect"
	"testing"
)

func PrimeFactor(n uint64) []uint64 {
	factors := make([]uint64, 0)
	for i := uint64(2); i <= n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
			for n%i == 0 {
				n /= i
			}
		}
	}
	return factors
}

func PrimitiveRoot(modulo uint64) uint64 {
	s := modulo - 1

	factors := PrimeFactor(s)

	for candidate := uint64(2); candidate < modulo; candidate++ {
		for _, factor := range factors {
			if ModExp(candidate, s/factor, modulo) == 1 {
				goto skip
			}
		}
		return candidate
	skip:
	}

	panic("No primitive root found")
}

func NthRootOfUnity(n, modulo uint64) uint64 {
	return ModExp(PrimitiveRoot(modulo), (modulo-1)/n, modulo)
}

func Range(n int) []uint64 {
	a := make([]uint64, n)
	for i := range a {
		a[i] = uint64(i)
	}
	return a
}

func TestNtt(t *testing.T) {
	N := 1024
	a := Range(N)
	// TODO(nsamar): Clean up arguments, they are messy
	modulo := NttFriendlyPrimes(N, 1, IntLog2(N)+4)[0]
	w := NthRootOfUnity(uint64(N), modulo)
	b := make([]uint64, len(a))
	copy(b, a)
	Ntt(a, w, modulo)
	TrivialNtt(b, w, modulo)
	if !reflect.DeepEqual(a, b) {
		t.Error("NTT failed\n", a, "\n", b)
	}
	if InverseTrivialNtt(b, w, modulo); !reflect.DeepEqual(b, Range(N)) {
		t.Error("Trivial InverseNTT failed\n", b, "\n", Range(N))
	}
	if InverseNtt(a, w, modulo); !reflect.DeepEqual(a, Range(N)) {
		t.Error("Inverse NTT failed", a, Range(N))
	}
}
