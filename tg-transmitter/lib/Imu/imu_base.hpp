#pragma once

#include "imu_data.hpp"

class ImuBase {
public:
    ImuBase() = default;

    virtual ImuData GetData() { return ImuData{}; }
};
