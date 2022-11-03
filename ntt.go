package fnttw

func Ntt(a []uint64, w uint64, modulo uint64) {
	n := len(a)
	if n == 1 {
		return
	}
	// Split the input in two halves
	half := n / 2
	a0 := a[:half]
	a1 := a[half:]
	// Compute the NTT of the two halves
	Ntt(a0, w*w%modulo, modulo)
	Ntt(a1, w*w%modulo, modulo)
	// Combine the two halves
	wi := uint64(1)
	for i := 0; i < half; i++ {
		t := wi * a1[i] % modulo
		a[i+half] = a[i] + modulo - t
		a[i] += t
		wi = wi * w % modulo
	}
}
