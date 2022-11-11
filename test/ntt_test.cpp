#include "gtest/gtest.h"
#include "ntt.h"

#include <random>

std::vector<uint32_t> RandomVector(int N) {
    std::vector<uint32_t> result;
    result.reserve(N);
    for (int i = 0; i < N; ++i) {
        result.push_back(std::rand());
    }
    return result;
}

void CheckEqual(const std::vector<uint32_t>& lhs, const std::vector<uint32_t>& rhs) {
  ASSERT_EQ(lhs.size(), rhs.size());
  for (int idx = 0; idx < lhs.size(); ++idx) {
    ASSERT_EQ(lhs[idx], rhs[idx]);
  }
}

constexpr uint32_t modulus = 1;

uint32_t IthNttTerm(const std::vector<uint32_t>& vec, uint32_t w, uint32_t i) {
	uint32_t result = 0;
	uint32_t wi = 1;
	w = ModExp<modulus>(w, i);
	for (const auto elem : vec) {
		result = ModAdd(result, ModMul<modulus>(elem, wi))
		wi = ModMul<modulus>(wi, w)
	}
	return result
}

TEST(NTT, NTT) {
  int N = 1 << 20;
  auto vec = RandomVector(N);
  auto ntt_vec = vec;
  std::vector<uint32_t> check;
  auto w = NthRootOfUnity<modulus>(modulus);
  for (int i = 0; i < N; ++i) {
      check.push_back(IthNttTerm(vec, w, modulus, i)
  }
  Ntt<modulus>(ntt_vec);
  CheckEqual(ntt_vec, check);
}
