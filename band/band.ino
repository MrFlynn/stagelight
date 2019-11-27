#include <Adafruit_NeoPixel.h>
#include <RHDatagram.h>
#include <RH_RF95.h>
#include <SPI.h>

// Setup LED strip.
#define LED_PIN 14
#define LED_COUNT 1

Adafruit_NeoPixel strip(LED_COUNT, LED_PIN, NEO_GRB + NEO_KHZ800);

// Setup packet radio.
#define __AVR_ATmega1284__

#define FREQUENCY 915.0
#define CS 4
#define INT 10

// Error display pin.
#define ERR 0

RH_RF95 rf95(CS, INT);
RHDatagram manager(rf95, 1);

// Commands for wireless communication.
enum Command {
  red,
  green,
  blue,
  vote
};

// Input buttons.
#define BUTTON1 25
#define BUTTON2 24

void display_error(uint8_t count) {
  while (1) {
    for (uint8_t i = 0; i < count; i++) {
      digitalWrite(ERR, HIGH);
      delay(500);
      digitalWrite(ERR, LOW);
      delay(500);
    }

    delay(500);
  }
}

// Button state machine.
enum ButtonStates { idle, hold1, hold2 } button_state;
volatile uint8_t last_vote = 0;

void voting_state() {
  // Transitions
  switch (button_state) {
    case idle:
    case hold1:
    case hold2:
    default:
      if (!digitalRead(BUTTON1)) {
        button_state = hold1;
      } else if (!digitalRead(BUTTON2)) {
        button_state = hold2;
      } else {
        button_state = idle;
      }

      break;
  }

  // Actions
  switch (button_state) {
    case hold1:
      last_vote = 1;
      digitalWrite(ERR, HIGH);
      break;
    case hold2:
      last_vote = 2;
      break;
    default:
      break;
  }
}

void setup() {
  strip.begin();

  // Error LED.
  pinMode(ERR, OUTPUT);

  // Setup button pins.
  pinMode(BUTTON1, INPUT_PULLUP);
  pinMode(BUTTON2, INPUT_PULLUP);

  if (!manager.init()) {
    display_error(1);
  }

  if (!rf95.setFrequency(FREQUENCY)) {
    display_error(2);
  }

  rf95.setTxPower(23, true);
}

void loop() {
  uint8_t buf[RH_RF95_MAX_MESSAGE_LEN];
  uint8_t len = sizeof(buf);
  uint8_t from;

  voting_state();
  
  if (manager.recvfrom(buf, &len, &from)) {
    if (!len) {
      return;
    }

    buf[len] = 0;

    if (buf[0] == red) {
      strip.setPixelColor(0, 255, 0, 0);
    } else if (buf[0] == green) {
      strip.setPixelColor(0, 0, 255, 0);
    } else if (buf[0] == blue) {
      strip.setPixelColor(0, 0, 0, 255);
    } else if (buf[0] == vote) {
      uint8_t data[] = {last_vote};

      manager.sendto(data, sizeof(data), from);
    }
  }
  
  strip.show();
}
