#include <RadioHandler.h>

#include <Scheduler.h>
#include <Task.h>
#include <LedTask.h>

// Radio information.
#define FREQUENCY 915.0
#define CS 4
#define INT 10
#define DEVICE_ID 1

// Task and scheduler setup.
Task **tasks;
Scheduler *scheduler;
RadioHandler *handler;

void setup() {  
  tasks = new Task*[1];
  tasks[0] = new LedTask;

  scheduler = new Scheduler(tasks, 1, 500);
  handler = new RadioHandler(FREQUENCY, CS, INT, DEVICE_ID);
}

void loop() {
  radioresult_t result = handler->receive();
  tasks[0]->importStream(result.buf, *result.len);
  
  scheduler->run();
  delay(scheduler->period());
}
