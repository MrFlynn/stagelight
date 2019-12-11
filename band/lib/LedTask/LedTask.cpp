#include "LedTask.h"

LedTask::LedTask(uint8_t pin, uint8_t count) : strip(Adafruit_NeoPixel(count, pin, NEO_GRB + NEO_KHZ800)),
    taskPeriod(500), 
    timeElapsed(0), 
    size(0),
    idx(0) {

    this->colors = new uint8_t*[42];
    for (int i = 0; i < 42; i++) {
        this->colors[i] = new uint8_t[3];
    }

    this->strip.begin();
    this->strip.setBrightness(64);
    this->strip.show();
}

void LedTask::nextTask() {
    // Transitions
    switch (this->state) {
        case INIT:
        case NEXT:
        default:
            this->state = NEXT;
            break;
    }

    // Actions
    switch (this->state) {
        case NEXT: {
            if (this->size == 0) {
                break;
            }

            uint8_t *color = this->colors[this->idx];
            this->strip.setPixelColor(0, color[0], color[1], color[2]);
            this->strip.show();

            this->idx++;
            if (this->idx >= this->size) {
                this->idx = 0;
            }
            
            }
            break;
        case INIT:
        default:
            break;
    }
}

void LedTask::importStream(uint8_t stream[], uint8_t len) {
    // Ignore the first two bytes are they are metadata.
    for (int i = 2; i < len; i += 3) {
        uint8_t *color = new uint8_t[3];
        color[0] = stream[i]; color[1] = stream[i + 1]; color[2] = stream[i + 2];
        
        this->colors[i / 3] = color;
    }

    this->size = (len - 1) / 3;
    this->idx = 0;
}

int LedTask::period() {
    return this->taskPeriod;
}

int LedTask::elapsed() {
    return this->timeElapsed;
}

void LedTask::setElapsed(int val) {
    this->timeElapsed = val;
}

uint8_t LedTask::getValue() {
    return 0;
}
