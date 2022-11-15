#include "ntt.h"
#include "utils.h"

constexpr int N = 1 << 22;
constexpr uint32_t modulus = 270532609;

int main() {
    std::vector<uint32_t> vec = RandomVector(N);
    while (true) {
        Ntt<modulus>(vec);
    }
}
