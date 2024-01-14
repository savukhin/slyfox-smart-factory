#pragma once

#ifndef ARDUINO
#include <ctime>
#define millis() std::time(nullptr)
#endif

#include "imu_data.hpp"
#include "mic_data.hpp"

#include <memory>

struct BufferItem {
    ImuData imu;
    MicData mic;
    uint64_t timestamp;
};

struct BufferItems {
    BufferItem *items_;
    size_t size;
};

class Buffer {
private:
    size_t size_ = 0;
    BufferItem *items_;
    size_t current_ = 0;

    BufferItems copyItems() {
        BufferItem *result = new BufferItem[current_];
        memcpy(result, items_, current_);
        return BufferItems{result, current_};
    }

public:
    Buffer(size_t size=100): size_(size) {
        recreate();
    }

    ~Buffer() {
        delete[] this->items_;
    }

    void recreate() {
        if (this->items_ == nullptr)
            delete[] this->items_;
        this->items_ = new BufferItem[size_];
    }

    BufferItems insert(ImuData &imu, MicData &mic) {
        BufferItem item{imu, mic};
        items_[current_++] = item;
        if (current_ < size_)
            return BufferItems{nullptr, 0};
        
        BufferItems copy = this->copyItems();
        recreate();
        return copy;
    }

};
