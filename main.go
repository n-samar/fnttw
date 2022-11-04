package main

import (
	"flag"
	"os"
	"runtime/pprof"

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
	profile_file, _ := os.Create("./cpu.pprof")
	pprof.StartCPUProfile(profile_file)
	fnttw.Ntt(a, modulus)
	pprof.StopCPUProfile()
}
