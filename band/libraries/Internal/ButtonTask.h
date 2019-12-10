#ifndef __BUTTON_TASK_H__
#define __BUTTON_TASK_H__

#include "Task.h"

class ButtonTask : public Task {
    private:
        enum states {
            INIT,
            WAIT,
            HOLDPOS,
            HOLDNEG
        } state;

        bool on;
        VoteValue value;

        int buttonPositive;
        int buttonNegative;

        int taskPeriod;
        int timeElapsed;
    public:
        enum VoteValue {
            NONE,
            POSITIVE,
            NEGATIVE
        };

        ButtonTask();
        void nextTask();
        void importStream(uint8_t[], uint8_t);
        int period();
        int elapsed();
        void setElapsed(int);
};

#endif
