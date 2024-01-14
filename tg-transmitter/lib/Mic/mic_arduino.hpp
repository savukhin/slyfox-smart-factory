#pragma once

#include "mic_base.hpp"

class MicArduino : public MicBase {
public:
    MicArduino() = default;

    MicData GetData() override { 
        return MicData{};
    }
};
