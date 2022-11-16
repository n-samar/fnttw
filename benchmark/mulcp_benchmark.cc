#include <benchmark/benchmark.h>

#include "utils.h"

typedef unsigned __int128 uint128_t;

constexpr int N = 1 << 20;

uint128_t a[N];
uint64_t b[N];
uint128_t c[N];

static void BM_MulCP(benchmark::State& state) {
    if (state.thread_index() == 0) {
        for (int i = 0; i < N; ++i) {
            a[i] = (uint128_t(rand()) << 64) + rand();
            b[i] = rand();
        }
    }
    for (auto _ : state) {
        for (int i = 0; i < N; ++i) {
            c[i] = a[i] * b[i];
        }
    }
}

BENCHMARK(BM_MulCP)->ThreadRange(1, 32)->Unit(benchmark::kMillisecond);
BENCHMARK_MAIN();
