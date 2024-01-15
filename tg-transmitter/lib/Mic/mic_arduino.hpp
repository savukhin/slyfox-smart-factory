#pragma once

#include <Arduino.h>
#include "mic_base.hpp"

class MicArduino : public MicBase {
private:
    uint8_t pin;
public:
    MicArduino(uint8_t pin) {
        this->pin = pin;
        pinMode(this->pin, INPUT);
    };

    MicData GetData() override { 
        int value = digitalRead(pin);
        return MicData{value};
    }
};
