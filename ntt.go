package fnttw

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

func NttHelper(a []uint64, w uint64, modulo uint64) {
	n := len(a)
	if n == 1 {
		return
	}
	// Split the input in two halves
	half := n / 2
	a0 := a[:half]
	a1 := a[half:]
	// Compute the NTT of the two halves
	NttHelper(a0, ModMul(w, w, modulo), modulo)
	NttHelper(a1, ModMul(w, w, modulo), modulo)
	// Combine the two halves
	wi := uint64(1)
	for i := 0; i < half; i++ {
		t := ModMul(wi, a1[i], modulo)
		a[i+half] = ModSub(a[i], t, modulo)
		a[i] = ModAdd(a[i], t, modulo)
		wi = ModMul(wi, w, modulo)
	}
}

func Ntt(a []uint64, w uint64, modulo uint64) {
	NttBitShuffle(a)
	NttHelper(a, w, modulo)
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

func TrivialNtt(a []uint64, w uint64, modulo uint64) {
	result := make([]uint64, len(a))
	for i := range a {
		result[i] = KthNttTerm(a, w, modulo, uint64(i))
	}
	copy(a, result)
}

func InverseTrivialNtt(a []uint64, w uint64, modulo uint64) {
	TrivialNtt(a, ModInv(w, modulo), modulo)

	// Divide by len(a)
	inv_n := ModInv(uint64(len(a)), modulo)
	for i := range a {
		a[i] = ModMul(a[i], inv_n, modulo)
	}
}

func InverseNtt(a []uint64, w uint64, modulo uint64) {
	Ntt(a, ModInv(w, modulo), modulo)

	// Divide by len(a)
	inv_n := ModInv(uint64(len(a)), modulo)
	for i := range a {
		a[i] = ModMul(a[i], inv_n, modulo)
	}
}
