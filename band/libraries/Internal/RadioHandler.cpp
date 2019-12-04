#include "RadioHandler.h"

RadioHandler::RadioHandler(int frequency, int cs, int interrupt, int id) : radio(cs, interrupt), 
    manager(radio, id), 
    deviceID(id) {
    this->manager.init();
    this->radio.setFrequency(frequency);
    this->radio.setTxPower(23);
}

radioresult_t RadioHandler::receive() {
    uint8_t buf[RH_RF95_MAX_MESSAGE_LEN];
    uint8_t len, from;
    
    radioresult_t result = { false, NOPCOMMAND, NOPMODE, buf, &len, &from };

    if (this->manager.recvfrom(buf, &len, &from)) {
        if (!len) {
            return result;
        }

        buf[len] = 0;

        result.hasData = true;
        result.command = (RadioCommand)buf[0];
        result.mode = (ModeCommand)buf[1];
    }

    return result;
}

void RadioHandler::send(radioresult_t result) {
    if (result.hasData) {
        manager.sendto(result.buf, *result.len, *result.from);
    }
}

void RadioHandler::getAutoDiscoverMessage(radioresult_t *result) {
    result->len = (uint8_t*)3;
    result->hasData = true;
    result->buf[2] = this->deviceID;
}
