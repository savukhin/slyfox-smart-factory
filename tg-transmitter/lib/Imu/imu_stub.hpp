#pragma once

#include <Arduino.h>
#include "imu_base.hpp"

class ImuStub : public ImuBase {
public:
    ImuStub() {
    };

    ImuData GetData() override {
        Vector3_float acc{
            1, 2, 3
        };
        Vector3_float gyro{
            10, 20, 30
        };
        return ImuData{acc, gyro};
    }
};
