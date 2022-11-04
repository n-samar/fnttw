package fnttw

import (
	"sync"
)

func BitReverse(x, n int) int {
	y := 0
	for i := 0; i < n; i++ {
		y = (y << 1) | (x & 1)
		x >>= 1
	}
	return y
}

func IntLog2(x int) int {
	result := 0
	for x > 1 {
		result += 1
		x >>= 1
	}
	return result
}

func NttBitShuffle(a []uint64) {
	b := make([]uint64, len(a))
	copy(b, a)
	for i := range a {
		b[i] = a[BitReverse(i, IntLog2(len(a)))]
	}
	copy(a, b)
}

func NttWithoutBitShuffle(a []uint64, w uint64, modulo Modulus64Bit) {
	n := len(a)
	if n == 1 {
		return
	}
	a0 := a[:n/2]
	a1 := a[n/2:]

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		NttWithoutBitShuffle(a0, modulo.ModMul(w, w), modulo)
	}()
	go func() {
		defer wg.Done()
		NttWithoutBitShuffle(a1, modulo.ModMul(w, w), modulo)
	}()

	wg.Wait()

	// Combine the two halves
	wi := modulo.FieldElement(1)
	for i := 0; i < n/2; i++ {
		t := modulo.ModMul(wi, a1[i])
		a[i+n/2] = ModSub(a[i], t, modulo.Modulus())
		a[i] = ModAdd(a[i], t, modulo.Modulus())
		wi = modulo.ModMul(wi, w)
	}
}

func NttTwiddle(a []uint64, w uint64, modulo Modulus64Bit) {
	NttBitShuffle(a)
	NttWithoutBitShuffle(a, w, modulo)
}

func Ntt(a []uint64, modulo Modulus64Bit) {
	NttTwiddle(a, modulo.NthRootOfUnity(uint64(len(a))), modulo)
}

func KthNttTerm(a []uint64, w uint64, modulo uint64, k uint64) uint64 {
	result := uint64(0)
	wi := uint64(1)
	w = ModExp(w, k, modulo)
	for i := range a {
		result = ModAdd(result, ModMul(a[i], wi, modulo), modulo)
		wi = ModMul(wi, w, modulo)
	}
	return result
}

func TrivialNttTwiddle(a []uint64, w uint64, modulo uint64) {
	result := make([]uint64, len(a))
	for i := range a {
		result[i] = KthNttTerm(a, w, modulo, uint64(i))
	}
	copy(a, result)
}

func TrivialNtt(a []uint64, modulo uint64) {
	w := ComputeNthRootOfUnity(uint64(len(a)), modulo)
	TrivialNttTwiddle(a, w, modulo)
}

func InverseTrivialNtt(a []uint64, modulo uint64) {
	w := ComputeNthRootOfUnity(uint64(len(a)), modulo)
	TrivialNttTwiddle(a, ModInv(w, modulo), modulo)

	// Divide by len(a)
	n_inverse := ModInv(uint64(len(a)), modulo)
	for i := range a {
		a[i] = ModMul(a[i], n_inverse, modulo)
	}
}

func InverseNtt(a []uint64, modulo Modulus64Bit) {
	w := modulo.NthRootOfUnity(uint64(len(a)))
	w_inv := modulo.ModInv(w)
	NttTwiddle(a, w_inv, modulo)

	// Divide by len(a)
	n_inverse := modulo.ModInv(modulo.FieldElement(uint64(len(a))))
	for i := range a {
		a[i] = modulo.ModMul(a[i], n_inverse)
	}
}
