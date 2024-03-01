#pragma once

#ifndef ARDUINO
#include <ctime>
#define millis() std::time(nullptr)
#endif

#include "imu_data.hpp"
#include "mic_data.hpp"

#include <cstring>
#include <utility>

struct BufferItem {
    ImuData imu;
    MicData mic;
    uint64_t timestamp;
};
template <typename _Tp, typename _Up = _Tp>
    inline _Tp
exchange(_Tp& __obj, _Up&& __new_val)
    { return std::__exchange(__obj, std::forward<_Up>(__new_val)); }

    
class BufferItems {
public:
    BufferItem *items_;
    size_t size;

    BufferItems(BufferItem *items, size_t size): items_(items), size(size) {
        // Serial.println("Calling constructor");
    }
    BufferItems(BufferItems& other): items_(other.items_), size(other.size) {
        // Serial.println("Calling copy constructor");
    }
    BufferItems(BufferItems&& other): items_(exchange(other.items_, nullptr)), size(exchange(other.size, 0)) {
        // Serial.println("Calling move constructor");
    }

    ~BufferItems() {
        // Serial.printf("Calling destructor %d\n", items_);
        // for (int i = 0; i < size; i++) {
        //     Serial.printf("%d ", items_[i]);
        // }
        // Serial.printf("\n");

        delete[] items_;
    }
};

class Buffer {
private:
    size_t size_ = 0;
    BufferItem *items_;
    size_t current_ = 0;

    BufferItems copyItems() {
        BufferItem *result = new BufferItem[current_];
        // Serial.printf(" memcpy");
        memcpy(result, items_, current_);
        // Serial.printf(" memcped");
        return BufferItems(result, current_);
    }

public:
    Buffer(size_t size=100): size_(size) {
        recreate();
    }

    ~Buffer() {
        delete[] this->items_;
    }

    void recreate() {
        if (this->items_ != nullptr)
            delete[] this->items_;
        this->items_ = new BufferItem[size_];
        this->current_ = 0;
    }

    size_t size() { return current_; }

     BufferItems insert(ImuData &imu, MicData &mic) {
        BufferItem item{imu, mic};
        items_[current_++] = item;
        if (current_ < size_)
            return BufferItems{nullptr, 0};
        // Serial.printf("copying ");
        BufferItems copy = this->copyItems();
        // Serial.printf(" recreating");
        recreate();
        // Serial.printf(" recreated");
        return copy;
    }

};
