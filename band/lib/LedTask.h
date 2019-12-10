#ifndef __LED_TASK_H__
#define __LED_TASK_H__

#define NUM_LEDS 1
#define LED_PIN 14

#include <Adafruit_NeoPixel.h>

#include "Task.h"

class LedTask : public Task {
    private:    
        enum states {
            INIT,
            NEXT
        } state;

        uint8_t **colors;
        uint8_t size;
        uint8_t idx;

        int taskPeriod;
        int timeElapsed;
    public:
        LedTask();
        void nextTask();
        void importStream(uint8_t[], uint8_t);
        int period();
        int elapsed();
        void setElapsed(int);
};

#endif
