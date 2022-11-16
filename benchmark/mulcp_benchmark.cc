#include <benchmark/benchmark.h>

#include "utils.h"

typedef unsigned __int128 uint128_t;

constexpr int N = 1 << 20;

uint128_t a[N];
uint64_t b[N];
uint128_t c[N];
uint32_t d[N];
uint32_t e[N][4];
uint32_t f[N][4];

static void BM_MulCP_128_64(benchmark::State& state) {
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

static void BM_MulCP_128_32(benchmark::State& state) {
    if (state.thread_index() == 0) {
        for (int i = 0; i < N; ++i) {
            a[i] = (uint128_t(rand()) << 64) + rand();
            d[i] = rand();
        }
    }
    for (auto _ : state) {
        for (int i = 0; i < N; ++i) {
            c[i] = a[i] * d[i];
        }
    }
}

static void BM_MulCP_4_32_32(benchmark::State& state) {
    if (state.thread_index() == 0) {
        for (int i = 0; i < N; ++i) {
            for (int j = 0; j < 4; ++j) {
                e[i][j] = rand();
            }
            d[i] = rand();
        }
    }
    for (auto _ : state) {
        for (int i = 0; i < N; ++i) {
            for (int j = 0; j < 4; ++j) {
                f[i][j] = a[i] * e[i][j];
            }
        }
    }
}

BENCHMARK(BM_MulCP_128_64)->ThreadRange(1, 32)->Unit(benchmark::kMillisecond);
BENCHMARK(BM_MulCP_128_32)->ThreadRange(1, 32)->Unit(benchmark::kMillisecond);
BENCHMARK(BM_MulCP_4_32_32)->ThreadRange(1, 32)->Unit(benchmark::kMillisecond);
BENCHMARK_MAIN();
