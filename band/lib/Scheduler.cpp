#include "Scheduler.h"

Scheduler::Scheduler(Task **tasks, int numTasks, int period) : tasks(tasks), 
    numTasks(numTasks),
    overallPeriod(period) {

}

void Scheduler::run() {
    for (int i = 0; i < this->numTasks; i++) {
        if (tasks[i]->period() >= tasks[i]->elapsed()) {
            tasks[i]->nextTask();
            tasks[i]->setElapsed(0);
        }

        tasks[i]->setElapsed(tasks[i]->elapsed() + overallPeriod);
    }
}

int Scheduler::period() {
    return this->overallPeriod;
}
