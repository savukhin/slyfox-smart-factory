#pragma once

#include <Arduino.h>
#include "imu_base.hpp"

#include <MPU9250.h>
// #include <Wire.h> 

class ImuArduino : public ImuBase {
private:
    MPU9250 device;
public:
    ImuArduino(): device(MPU9250(Wire, 0x68)) {
    };

    ImuData GetData() override {
        Vector3_float acc{
            device.getAccelX_mss(),
            device.getAccelY_mss(),
            device.getAccelZ_mss()
        };
        Vector3_float gyro{
            device.getGyroX_rads(),
            device.getGyroX_rads(),
            device.getGyroX_rads()
        };
        return ImuData{acc, gyro};
    }
};
