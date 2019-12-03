#ifndef __RADIO_HANDLER_H__
#define __RADIO_HANDLER_H__

#include <RHDatagram.h>
#include <RH_RF95.h>
#include <SPI.h>

// Setup packet radio.
#define __AVR_ATmega1284__

enum RadioCommand {
    NOPCOMMAND,
    DISCOVER,
    STATUSUPDATE
};

enum ModeCommand {
    NOPMODE,
    CHANGECOLOR,
    VOTESTATUS
};

struct radioresult_t {
    bool hasData;
    RadioCommand command;
    ModeCommand mode;
    uint8_t *buf, *len, *from;
};
    

class RadioHandler {
    private:
        RH_RF95 radio;
        RHDatagram manager;

        uint8_t deviceID;
    public:
        RadioHandler(int, int, int, int);
        radioresult_t receive();
        void send(radioresult_t);
        void getAutoDiscoverMessage(radioresult_t*);
};

#endif
