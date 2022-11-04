package fnttw

func IsPrime(n uint64) bool {
	return NewBigIntFromUint64(n).ProbablyPrime(20)
}

func NttFriendlyPrimes(ntt_size, num_primes, bits_per_prime int) []uint64 {
	primes := make([]uint64, num_primes)
	if bits_per_prime > 63 {
		panic("bits_per_prime must be <= 63")
	}
	m := (1 << (bits_per_prime - 1)) / ntt_size
	for i := 0; i < num_primes; i++ {
		candidate := uint64(0)
		for !IsPrime(candidate) {
			candidate = uint64(2*ntt_size*m + 1)
			m += 1
		}
		primes[i] = candidate
	}
	return primes
}
