#include "c_radio.h"
#include "RadioHead/RH_RF69.h"

C_RADIO c_radio(uint8_t slaveSelectPin, uint8_t interruptPin) {
    RH_RF69 *ret = new RH_RF69(slaveSelectPin, interruptPin);
    return (void*)ret;
}