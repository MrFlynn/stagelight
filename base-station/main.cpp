#include <thread>
#include <chrono>

#include <bcm2835.h>
#include <RH_RF69.h>

#define FREQUENCY 915.0

RH_RF69 rf69(BCM2835_SPI_CS0, RPI_V2_GPIO_P5_06);

int main(void) {
    uint8_t * data = 1;

    if (!rf69.init()) {
        std::cout << "Radio failed to initialize" << std::endl;
        return 1;
    }

    if (!rf69.setFrequency(FREQUENCY)) {
        std::cout << "Failed to set frequency of radio" << std::endl;
        return 1;
    }

    while (true) {
        rf69.send(data, sizeof(uint8_t));
        rf69.waitPacketSent();

        std::this_thread::sleep_for(std::chrono::seconds(1));
    }    
}