#ifndef __LED_TASK_H__
#define __LED_TASK_H__

#include <Adafruit_NeoPixel.h>

#include <Task.h>

class LedTask : public Task {
    private:    
        enum states {
            INIT,
            NEXT
        } state;

        Adafruit_NeoPixel strip;
        
        int taskPeriod;
        int timeElapsed;

        uint8_t **colors;
        uint8_t size;
        uint8_t idx;
    public:
        LedTask(uint8_t, uint8_t);
        void nextTask();
        void importStream(uint8_t[], uint8_t);
        int period();
        int elapsed();
        void setElapsed(int);
        uint8_t getValue();
};

#endif
