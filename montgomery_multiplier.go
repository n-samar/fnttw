package fnttw

import (
	"math/big"
	"math/bits"
)

type Montgomery64BitModulus struct {
	q         uint64
	q_inverse uint64
}

type Simple64BitModulus struct {
	q uint64
}

type Modulus64Bit interface {
	ModMul(a, b uint64) uint64
	Modulus() uint64
	NthRootOfUnity(n uint64) uint64
	ModInv(a uint64) uint64
	FieldElement(a uint64) uint64
	FromFieldElement(a uint64) uint64
}

func (modulus Simple64BitModulus) FieldElement(a uint64) uint64 {
	return a
}

func (modulus Simple64BitModulus) FromFieldElement(a uint64) uint64 {
	return a
}

func (modulus Montgomery64BitModulus) FromFieldElement(a uint64) uint64 {
	return modulus.ModMul(a, uint64(1))
}

func (modulus Montgomery64BitModulus) FieldElement(a uint64) uint64 {
	return ModMul(a, modulus.RMod(), modulus.q)
}

func (modulus Simple64BitModulus) NthRootOfUnity(n uint64) uint64 {
	return ComputeNthRootOfUnity(n, modulus.q)
}

func (modulus Simple64BitModulus) ModInv(a uint64) uint64 {
	return ModInv(a, modulus.q)
}

func (modulus Montgomery64BitModulus) ModInv(a uint64) uint64 {
	a_non_montgomery := modulus.FromFieldElement(a)
	a_inverse := ModInv(a_non_montgomery, modulus.q)
	return modulus.FieldElement(a_inverse)
}

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

func ComputeNthRootOfUnity(n, modulo uint64) uint64 {
	return ModExp(PrimitiveRoot(modulo), (modulo-1)/n, modulo)
}

func (modulus Montgomery64BitModulus) NthRootOfUnity(n uint64) uint64 {
	return modulus.FieldElement(ComputeNthRootOfUnity(n, modulus.q))
}

func (modulus Simple64BitModulus) Modulus() uint64 {
	return modulus.q
}

func (modulus Montgomery64BitModulus) Modulus() uint64 {
	return modulus.q
}

func R64Bit() *big.Int {
	return big.NewInt(0).Mul(big.NewInt(1<<62), big.NewInt(1<<2))
}

func (modulus Montgomery64BitModulus) R() *big.Int {
	return R64Bit()
}

func (modulus Montgomery64BitModulus) RMod() uint64 {
	return big.NewInt(0).Mod(modulus.R(), NewBigIntFromUint64(modulus.q)).Uint64()
}

func MakeMontgomery64BitModulus(q uint64) Montgomery64BitModulus {
	R := R64Bit()
	R_inverse := MultiplicativeInverse(NewBigIntFromUint64(q), R)

	return Montgomery64BitModulus{q, R_inverse}
}

func NewBigIntFromUint64(a uint64) *big.Int {
	return big.NewInt(0).SetUint64(a)
}

func MultiplicativeInverse(a, modulo *big.Int) uint64 {
	big_result := big.NewInt(0)
	return big_result.ModInverse(a, modulo).Uint64()
}

func ModMul(a, b, q uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return bits.Rem64(hi, lo, q)
}

func ModAdd(a, b, q uint64) uint64 {
	return (a + b) % q
}

func ModSub(a, b, q uint64) uint64 {
	if b > a {
		return q - (b - a)
	}
	return a - b
}

func ModInv(a, modulo uint64) uint64 {
	big_a := NewBigIntFromUint64(a)
	big_modulo := NewBigIntFromUint64(modulo)
	return big_a.ModInverse(big_a, big_modulo).Uint64()
}

func ModExp(a, n, q uint64) uint64 {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return a % q
	}
	sqrt := ModExp(a, n/2, q)
	result := ModMul(sqrt, sqrt, q)
	if n%2 == 1 {
		result = ModMul(result, a, q)
	}
	return result
}

func (modulus Simple64BitModulus) ModMul(a, b uint64) uint64 {
	return ModMul(a, b, modulus.q)
}

func (modulus Montgomery64BitModulus) ModMul(a, b uint64) uint64 {
	// From https://github.com/tuneinsight/lattigo/blob/master/ring/modular_reduction.go#L60
	mhi, mlo := bits.Mul64(a, b)
	hhi, _ := bits.Mul64(mlo*modulus.q_inverse, modulus.q)
	result := mhi - hhi + modulus.q
	if result >= modulus.q {
		result -= modulus.q
	}
	return result
}
