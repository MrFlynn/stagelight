#include "ButtonTask.h"

ButtonTask::ButtonTask(int buttonPositive, int buttonNegative) : on(false),
    buttonPositive(buttonPositive),
    buttonNegative(buttonNegative),
    value(0),
    taskPeriod(50),
    elapsedTime(0) {

    pinMode(buttonPositive, INPUT_PULLUP);
    pinMode(buttonNegative, INPUT_PULLUP);
}

void ButtonTask::nextTask() {
    // Transitions
    switch (this->state) {
        case WAIT:
            if (!digitalRead(buttonPositive)) {
                this->state = HOLDPOS;
            } else if (!digitalRead(buttonNegative)) {
                this->state = HOLDNEG;
            } else {
                this->state = WAIT;
            
            break;
        case HOLDPOS:
            if digitalRead(buttonPositive) {
                this->state = WAIT;
            } else {
                this->state = HOLDPOS;
            }

            break;
        case HOLDNEG:
            if digitalWrite(buttonNegative) {
                this->state = WAIT;
            } else {
                this->state = HOLDNEG;
            }

            break;
        case INIT:
        default:
            this->state = WAIT;
            break;
    }

    // Actions:
    switch (this->state) {
        case WAIT:
            this->value = NONE;
            break;
        case HOLDPOS:
            this->value = POSITIVE;
            break;
        case HOLDNEG:
            this->value = NEGATIVE;
            break;
        default:
            this->value = NONE;
            break;
    }
}

void ButtonTask::importStream(uint8_t[] stream, uint8_t len) {
    return;
}

int ButtonTask::period() {
    return this->taskPeriod;
}

int ButtonTask::elapsed() {
    return this->elapsedTime;
}

void ButtonTask::setElapsed(int val) {
    this->elapsedTime = val;
}

uint8_t ButtonTask::getValue() {
    if (this->on) {
        return (uint8_t)this->value;
    }

    return 0;
}
