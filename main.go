package main

import (
	"flag"

	"github.com/n-samar/fnttw/fnttw"
)

func main() {
	logN := flag.Int("logN", 10, "log number of NTT element")
	flag.Parse()
	N := 1 << *logN
	a := fnttw.Range(N)
	// TODO(nsamar): Clean up arguments, they are messy
	q := fnttw.NttFriendlyPrimes(N, 1, *logN+4)[0]
	modulus := fnttw.MakeMontgomery64BitModulus(q)
	fnttw.Ntt(a, modulus)
}
