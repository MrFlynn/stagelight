#ifndef __TASK_H__
#define __TASK_H__

#include <Arduino.h>

class Task {
    public:
        virtual void nextTask() = 0;
        virtual void importStream(uint8_t[], uint8_t) = 0;
        virtual int period() = 0;
        virtual int elapsed() = 0;
        virtual void setElapsed(int) = 0;
};

#endif
