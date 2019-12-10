#ifndef __SCHEDULER_H__
#define __SCHEDULER_H__

#include "Task.h"

class Scheduler {
    private:
        Task** tasks;
        int numTasks;
        int overallPeriod;

    public:
        Scheduler(Task**, int, int);

        int period();
        void run();
};

#endif
