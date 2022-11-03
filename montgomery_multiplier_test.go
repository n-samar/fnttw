package fnttw

import (
	"math/rand"
	"testing"
)



func TestMontgomeryModMul64(t *testing.T) {
	modulus := uint64(13)
	montgomery_modulus := MakeMontgomery64BitModulus(modulus)
	simple_modulus := Simple64BitModulus{modulus}

	for i := 0; i < 100; i++ {
		a := rand.Uint64() % modulus
		b := rand.Uint64() % modulus
		montgomery_output := montgomery_modulus.ModMul(a, b)
		simple_output := simple_modulus.ModMul(a, b)

		if simple_modulus.ModMul(montgomery_modulus.RMod(), montgomery_output) != simple_output {
			t.Fatalf("(%d) MontgomeryMultiplier(%d, %d) = %d, want %d", montgomery_modulus.q, a, b, montgomery_output, simple_output)
		}
	}
}
