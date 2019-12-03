#include <Scheduler.h>
#include <Task.h>
#include <LedTask.h>

Task **tasks;
Scheduler *scheduler;

void setup() {  
  tasks = new Task*[1];
  tasks[0] = new LedTask;

  scheduler = new Scheduler(tasks, 1, 500);
  
  uint8_t colors[10] = {0, 255, 0, 0, 0, 255, 0, 0, 0, 255};
  tasks[0]->importStream(colors, 10);
}

void loop() {
  scheduler->run();
  delay(scheduler->period());
}
