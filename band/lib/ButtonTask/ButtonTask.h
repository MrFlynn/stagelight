#ifndef __BUTTON_TASK_H__
#define __BUTTON_TASK_H__

#include <Task.h>

class ButtonTask : public Task {
    public:
        enum VoteValue {
            NONE,
            POSITIVE,
            NEGATIVE
        };
    private:
        enum states {
            INIT,
            WAIT,
            HOLDPOS,
            HOLDNEG
        } state;

        int buttonPositive;
        int buttonNegative;

        bool on;
        VoteValue value;

        int taskPeriod;
        int timeElapsed;
    public:
        ButtonTask(int, int);
        void nextTask();
        void importStream(uint8_t[], uint8_t);
        int period();
        int elapsed();
        void setElapsed(int);
        uint8_t getValue();
};

#endif
