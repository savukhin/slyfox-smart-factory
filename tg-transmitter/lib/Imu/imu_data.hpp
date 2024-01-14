#pragma once

#include <stdint.h>

struct Vector3_64 {
    int64_t x, y, z;
};

struct ImuData {
    Vector3_64 acc;
    Vector3_64 gyro;
};
