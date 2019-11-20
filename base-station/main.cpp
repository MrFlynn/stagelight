#include <thread>
#include <chrono>
#include <iostream>

#include <RH_RF69.h>

extern "C" {
#include <wiringPi.h>
}

#define FREQUENCY 915.0

#define CS 10
#define INT 6
#define RST 5

RH_RF69 rf69(CS, INT);

void reset_radio(void) {
    pinMode(RST, OUTPUT);

    digitalWrite(RST, HIGH);
    delay(10);
    digitalWrite(RST, LOW);
    delay(10);
}

int main(void) {
    RasPiSetup();

    reset_radio();

    if (!rf69.init()) {
        std::cout << "Radio failed to initialize." << std::endl;
        return 1;
    }

    if (!rf69.setFrequency(FREQUENCY)) {
        std::cout << "Failed to set frequency of radio." << std::endl;
        return 1;
    }

    while (true) {
        std::cout << "Sending packets" << std::endl;
        rf69.send((uint8_t *)1, sizeof(uint8_t));
        rf69.waitPacketSent();

        std::this_thread::sleep_for(std::chrono::seconds(1));
    }    
}
