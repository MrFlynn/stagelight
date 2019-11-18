#include "led.hpp"

void setup_pwm(void) {
    cli();

    // Enable OC2A/B ports (PD6 and PD7) and clear ports.
    DDRD |= _BV(PD6) | _BV(PD7); PORTD |= ~(_BV(PD6) | _BV(PD7));

    // Set comparator reset on match.
    // Enable fast PWM mode 3.
    TCCR2A = _BV(COM2A1) | _BV(COM2B1) | _BV(WGM21) | _BV(WGM20);

    // Enable fast PWM to cap at value of OCR2A.
    // Set prescaler to 1 (8MHz PWM).
    TCCR2B = _BV(WGM22) | _BV(CS20);

    // Enable interrupt on match of OCR2B.
    TIMSK2 = _BV(OCIEB);

    // Set MAX count value for PWM.
    OCR2A = 20;

    // Set default duty cycle to 50%.
    0CR2B = 10;

    sei();
}
