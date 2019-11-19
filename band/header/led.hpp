#include <avr/io.h>
#include <avr/interrupt.h>

extern "C" {
    #include <util/delay.h>
}

#ifndef LIBLED
#define LIBLED

#define T1H 9
#define T0H 4

class Led {
    private:
        // Array containing timing values for PWM interrupt.
        uint8_t pwm_times[44];
        uint8_t index;

        // Methods used internally to start and stop PWM.
        void start_pwm();
        void stop_pwm();

        // This method is used to fill out specific section of the pwm_times array.
        void calculate_pwm_times(uint8_t, uint8_t);
    public:
        // Constructor.
        Led();

        // Class methods.
        void set_color(uint8_t, uint8_t, uint8_t);
        void service_isr();
};

#endif
