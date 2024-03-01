#pragma once

#include <Arduino.h>
#include "mic_base.hpp"

class MicStub : public MicBase {
public:
    MicStub() {
    };

    MicData GetData() override { 
        return MicData{100};
    }
};
