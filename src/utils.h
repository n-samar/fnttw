#ifndef FNTTW_UTILS_H_
#define FNTTW_UTILS_H_

#include <vector>

inline std::vector<uint32_t> RandomVector(int N) {
    std::vector<uint32_t> result;
    result.reserve(N);
    for (int i = 0; i < N; ++i) {
        result.push_back(std::rand());
    }
    return result;
}

#endif  // FNTTW_UTILS_H_
