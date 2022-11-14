#include <benchmark/benchmark.h>
#include <omp.h>

#include "ntt.h"
#include "utils.h"

constexpr uint32_t modulus = 270532609;

std::vector<uint32_t> vec;

static void BM_Ntt(benchmark::State& state) {
    if (state.thread_index == 0) {
        vec = RandomVector(state.range(0));
    }
    for (auto _ : state) {
        Ntt<modulus>(vec);
    }
}

BENCHMARK(BM_Ntt)->Range(1 << 20, 1 << 25)->ThreadRange(1, 32)->Unit(benchmark::kMillisecond);
BENCHMARK_MAIN();
