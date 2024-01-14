#pragma once

#include "mic_data.hpp"

class MicBase {
public:
    MicBase() = default;

    virtual MicData GetData() { return MicData{}; }
};
