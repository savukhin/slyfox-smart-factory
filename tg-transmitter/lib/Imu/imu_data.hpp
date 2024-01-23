#pragma once

#include <stdint.h>

struct Vector3_64 {
    int64_t x, y, z;
};

struct Vector3_float {
    float x, y, z;
};

struct ImuData {
    Vector3_float acc;
    Vector3_float gyro;
};
