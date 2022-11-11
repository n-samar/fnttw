#include <benchmark/benchmark.h>

#include "ntt.h"
#include "utils.h"

constexpr uint32_t N = 1 << 20;
constexpr uint32_t modulus = 270532609;

static void BM_Ntt(benchmark::State& state) {
    auto vec = RandomVector(N);
    for (auto _ : state) {
        Ntt<modulus>(vec);
    }
}

BENCHMARK(BM_Ntt);
BENCHMARK_MAIN();
