package fnttw

import (
	"math/big"
	"math/bits"
)

type Montgomery64BitModulus struct {
	q         uint64
	R_inverse uint64
}

type Simple64BitModulus struct {
	q uint64
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

func (modulus Simple64BitModulus) ModMul(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return bits.Rem64(hi, lo, modulus.q)
}

func (modulus Montgomery64BitModulus) ModMul(a, b uint64) uint64 {
	// From https://github.com/tuneinsight/lattigo/blob/master/ring/modular_reduction.go#L60
	mhi, mlo := bits.Mul64(a, b)
	hhi, _ := bits.Mul64(mlo*modulus.R_inverse, modulus.q)
	result := mhi - hhi + modulus.q
	if result >= modulus.q {
		result -= modulus.q
	}
	return result
}
