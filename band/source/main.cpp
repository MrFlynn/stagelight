#include <led.hpp>

Led led = Led();

ISR(TIMER2_COMPA_vect) {
    led.service_isr();
}

int main(void) {
    led.set_color(0, 0, 255);
    while (true) {}

    return 1;
}