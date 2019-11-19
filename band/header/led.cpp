#include "led.hpp"

#define F_CPU 8000000UL

Led::Led() : pwm_times {0}, index(0) {
    cli();

    // Enable OC2A/B ports (PD6 and PD7) and clear ports.
    // DDRD |= _BV(PD6) | _BV(PD7); PORTD |= ~(_BV(PD6) | _BV(PD7));
    DDRD = 0xFF; PORTD = 0x00;

    // Set MAX count value for PWM.
    OCR2A = 20;

    // Set default duty cycle to 0%.
    OCR2B = 0;

    sei();
}

void Led::start_pwm() {
    // Set comparator reset on match.
    // Enable fast PWM mode 3.
    TCCR2A = _BV(COM2A1) | _BV(COM2B1) | _BV(WGM21) | _BV(WGM20);

    // Enable fast PWM to cap at value of OCR2A.
    // Set prescaler to 1 (8MHz PWM).
    TCCR2B = _BV(WGM22) | _BV(CS20);

    // Enable interrupt on match of OCR2B.
    TIMSK2 = _BV(OCIE2A);
    TIFR2 = _BV(OCF2A);
}

void Led::stop_pwm() {
    // Zero out all PWM registers to stop clock.
    TCCR2A = 0;
    TCCR2B = 0;
    TIMSK2 = 0;
}

void Led::calculate_pwm_times(uint8_t color, uint8_t offset) {
    for (uint8_t i = 0; i < 9; i++) {
        if (color & (0xFF >> i)) {
            this->pwm_times[offset + i] = T1H;
        } else {
            this->pwm_times[offset + i] = T0H;
        }
    }
}

void Led::set_color(uint8_t red, uint8_t green, uint8_t blue) {
    this->calculate_pwm_times(green, 0);
    this->calculate_pwm_times(red, 8);
    this->calculate_pwm_times(blue, 16);

    this->start_pwm();
}

void Led::service_isr() {
    // Set output compare register 2 to next calculated time.
    OCR2B = this->pwm_times[this->index];
    this->index++;

    if (this->index > 43) {
        this->index = 0;

        this->stop_pwm();
    }
}