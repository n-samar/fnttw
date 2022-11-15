#include "ntt.h"
#include "utils.h"

constexpr int N = 1 << 20;
constexpr uint32_t modulus = 270532609;

int main() {
    std::vector<uint32_t> vec = RandomVector(N);
    Ntt<modulus>(vec);
}
