#include <Arduino.h>

// Project includes.
#include <Task.h>
#include <Scheduler.h>
#include <RadioHandler.h>

// Specific tasks.
#include <LedTask.h>
#include <ButtonTask.h>

// Pins.
#define LED_PIN 7
#define LED_COUNT 2

#define BUTTON_POS 12
#define BUTTON_NEG 10

#define RADIO_CS 16
#define RADIO_INT 5

// Other defintions
#define FREQ 915.0
#define DEVICE_ID 1

// Define tasks struct.
Task **tasks;
Scheduler *scheduler;
RadioHandler radio(FREQ, RADIO_CS, RADIO_INT, DEVICE_ID);

radioresult_t handle_radio_command() {
  radioresult_t result = radio.receive();

  if (result.command == 1) {
    radio.getAutoDiscoverMessage(&result);
  } else if (result.command == 2) {
    if (result.mode == 1) {
      tasks[0]->importStream(result.buf, *result.len);
    } else if (result.mode == 2 || result.mode == 3) {
      tasks[1]->importStream(result.buf, *result.len);
      
      uint8_t newbuf[3] = {2, result.mode, tasks[1]->getValue()};
      result.buf = newbuf;
      result.len = new uint8_t(3);
    }
  }

  return result;
}

void setup() {
  // Setup tasks and scheduler.
  tasks = new Task*[2];
  tasks[0] = new LedTask(LED_PIN, LED_COUNT);
  tasks[1] = new ButtonTask(BUTTON_POS, BUTTON_NEG);
  scheduler = new Scheduler(tasks, 2, 50);
}

void loop() {
  // Get result and run scheduler.
  radioresult_t res = handle_radio_command();
  scheduler->run();

  // Send back data if the incoming command expects something back.
  if (res.command == 1 || res.mode == 2 || res.mode == 3) {
    radio.send(res);
  }

  // Delay for X ms until next scheduler execution.
  delay(scheduler->period());
}