#pragma once

#ifndef ARDUINO
#include <ctime>
#define millis() std::time(nullptr)
#endif

#include "imu_data.hpp"
#include "mic_data.hpp"

#include <inttypes.h>
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
        // Serial.printf("free space %d and sizeof new is %d", heap_caps_get_free_size(MALLOC_CAP_DEFAULT), sizeof(BufferItem) * current_);
        BufferItem *result = new BufferItem[current_];
        // Serial.printf("free space after new %d", heap_caps_get_free_size(MALLOC_CAP_DEFAULT));
        // Serial.printf(" memcpy");
        // Serial.printf("mic befor memcpy %" PRId64 " %" PRId64 " %d\n", result[0].mic.value, items_[0].mic.value, current_);
        // Serial.printf("mic befor memcpy %" PRId64 " %" PRId64 "\n", result[current_-1].mic.value, items_[current_-1].mic.value);
        memcpy(result, items_, sizeof(BufferItem) * current_);
        // Serial.printf("mic after memcpy %" PRId64 " %" PRId64 "\n", result[0].mic.value, items_[0].mic.value);
        // Serial.printf("mic after memcpy %" PRId64 " %" PRId64 "\n", result[current_-1].mic.value, items_[current_-1].mic.value);

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
        // Serial.printf("Before delete free space %d and sizeof new is %d\n", heap_caps_get_free_size(MALLOC_CAP_DEFAULT), sizeof(BufferItem) * size_);
        // if (this->items_ != nullptr)
        //     delete[] this->items_;

        // Serial.printf("Recreating free space %d and sizeof new is %d\n", heap_caps_get_free_size(MALLOC_CAP_DEFAULT), sizeof(BufferItem) * size_);
        // this->items_ = new BufferItem[size_];
        if (this->items_ == nullptr)
            this->items_ = new BufferItem[size_];
        // Serial.println("recreated");
        this->current_ = 0;
    }

    size_t size() { return current_; }
    size_t maxsize() { return size_; }

    BufferItems insert(ImuData &imu, MicData &mic) {
        BufferItem item{imu, mic};
        if (current_ < size_)
            items_[current_++] = item;
        // if (current_ < size_)
        //     return BufferItems{nullptr, 0};
        // Serial.printf("copying ");
        // BufferItems copy = this->copyItems();
        // Serial.printf(" recreating");
        // Serial.printf("mic after copy %" PRId64 "\n", copy.items_[current_-1].mic.value);

        // recreate();
        // Serial.printf(" recreated");
        // return copy;
        return BufferItems{nullptr, 0};
    }

    BufferItem get(size_t ind) {
        return items_[ind];
    }

    BufferItem getLast() {
        return this->items_[current_];
    }

};
