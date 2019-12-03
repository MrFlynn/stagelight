#include "LedTask.h"

Adafruit_NeoPixel strip(NUM_LEDS, LED_PIN, NEO_GRB + NEO_KHZ800);

LedTask::LedTask() : taskPeriod(500), size(0), idx(0), timeElapsed(0) {
    this->colors = new uint8_t*[42];
    for (int i = 0; i < 42; i++) {
        this->colors[i] = new uint8_t[3];
    }

    strip.begin();
    
    strip.setBrightness(64);
    strip.show();
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
            strip.setPixelColor(0, color[0], color[1], color[2]);
            strip.show();

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
    for (int i = 1; i < len; i += 3) {
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
