#ifndef FNTTW_NTT_H_
#define FNTTW_NTT_H_

#include <vector>
#include <iostream>
#include <stdint.h>
#include <omp.h>
#include <cassert>

#include <glog/logging.h>

inline std::vector<uint32_t> PrimeFactors(uint32_t n) {
    std::vector<uint32_t> factors;
    for (uint32_t i = 2; i <= n; ++i) {
        if (n % i == 0) {
            factors.push_back(i);
            while (n%i == 0) {
                n /= i;
            }
        }
    }
    return factors;
}

template <uint32_t modulus>
static inline uint32_t ModAdd(uint32_t a, uint32_t b) {
    return (a+b) % modulus;
}

template <uint32_t modulus>
static inline uint32_t ModSub(uint32_t a, uint32_t b) {
    if (b > a) {
        return modulus - (b-a);
    }
    return (a-b) % modulus;
}

template <uint32_t modulus>
static inline uint32_t ModMul(uint32_t a, uint32_t b) {
    uint64_t prod = a * uint64_t(b);
    return prod % modulus;
}

template <uint32_t modulus>
uint32_t ModExp(uint32_t a, uint32_t pow) {
	if (pow == 0) {
		return 1;
	}
	if (pow == 1) {
		return a % modulus;
	}
	auto sqrt = ModExp<modulus>(a, pow/2);
	auto result = ModMul<modulus>(sqrt, sqrt);
	if (pow%2 == 1) {
		result = ModMul<modulus>(result, a);
	}
	return result;
}

template <uint32_t modulus>
uint32_t PrimitiveRoot() {
    auto s = modulus - 1;
    auto factors = PrimeFactors(s);

    for (uint32_t candidate = 2; candidate < modulus; ++candidate) {
        for (const auto& factor : factors) {
            if (ModExp<modulus>(candidate, s / factor) == 1) {
                goto skip;
            }
        }
        return candidate;
skip:
        (void)0;
    }
    LOG(FATAL) << "No primitive root of modulus " << modulus;
}

template <uint32_t modulus>
uint32_t NthRootOfUnity(uint32_t N) {
    return ModExp<modulus>(PrimitiveRoot<modulus>(), (modulus-1)/N);
}

constexpr int elements_per_thread = 1 << 23;

template <uint32_t modulus>
void NttWithoutBitShuffle(uint32_t* vec, uint32_t n, uint32_t w) {
	if (n <= elements_per_thread) {
        for (int n_sub = 2; n_sub <= n; n_sub *= 2) {
            auto w_curr = ModExp<modulus>(w, n/n_sub);
            auto w_offset = ModExp<modulus>(w_curr, n_sub);
            auto w_base = uint32_t(1);
            for (int curr_offset = 0; curr_offset < n; curr_offset+=n_sub) {
                auto wi = w_base;
                for (int i = 0; i < n_sub/2; ++i) {
                    auto t = ModMul<modulus>(wi, vec[curr_offset + n_sub/2 + i]);
                    vec[curr_offset + i + n_sub/2] = ModSub<modulus>(vec[curr_offset + i], t);
                    vec[curr_offset + i] = ModAdd<modulus>(vec[curr_offset + i], t);
                    wi = ModMul<modulus>(wi, w_curr);
                }
                w_base = ModMul<modulus>(w_base, w_offset);
            }
        }
        return;
	}

    auto w_squared = ModMul<modulus>(w, w);
    #pragma omp parallel sections
    {
        #pragma omp section
        NttWithoutBitShuffle<modulus>(vec, n/2, w_squared);
        #pragma omp section
        NttWithoutBitShuffle<modulus>(vec + n/2, n/2, w_squared);
    }

    int num_threads = int(n/2/elements_per_thread);

    #pragma omp parallel for schedule(static, 1)
    for (int j = 0; j < num_threads; ++j) {
        int start_idx = j*elements_per_thread;
        auto wi = ModExp<modulus>(w, start_idx);
        for (int i = start_idx; i < start_idx + elements_per_thread; ++i) {
            auto t = ModMul<modulus>(wi, vec[n/2 + i]);
            vec[i+n/2] = ModSub<modulus>(vec[i], t);
            vec[i] = ModAdd<modulus>(vec[i], t);
            wi = ModMul<modulus>(wi, w);
        }
    }
}

inline uint32_t BitReverse(uint32_t bitwidth, uint32_t x) {
	uint32_t y = 0;
	for (uint32_t i = 0; i < bitwidth; i++) {
		y = (y << 1) | (x & 1);
		x >>= 1;
	}
	return y;
}

inline uint32_t log2(uint32_t n) {
    uint32_t result = 0;
    while (n >>= 1) ++result;
    return result;
}

inline void NttBitShuffle(std::vector<uint32_t>& vec) {
    std::vector<uint32_t> result;
    result.reserve(vec.size());
    for (int i = 0; i < vec.size(); ++i) {
        result.push_back(vec[BitReverse(log2(vec.size()), i)]);
     }
     vec = result;
}

template <uint32_t modulus>
void NttTwiddle(std::vector<uint32_t>& vec, uint32_t w) {
    NttBitShuffle(vec);
    NttWithoutBitShuffle<modulus>(&vec[0], vec.size(), w);
}

template <uint32_t modulus>
void Ntt(std::vector<uint32_t>& vec) {
    NttTwiddle<modulus>(vec, NthRootOfUnity<modulus>(vec.size()));
}


#endif  // FNTTW_NTT_H_
