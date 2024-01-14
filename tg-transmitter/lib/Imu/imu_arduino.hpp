#pragma once

#include "imu_base.hpp"

class ImuArduino : public ImuBase {
public:
    ImuArduino() = default;

    ImuData GetData() override { 
        return ImuData{};
    }
};
