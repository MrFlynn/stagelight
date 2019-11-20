RH 		= ../include/RadioHead

CC		= g++
CFLAGS 	= -DRASPBERRY_PI -DBCM2835_NO_DELAY_COMPATIBILITY -Wall --std=c++11 
INCLUDE = -I$(RH)
LIBS	= -lwiringPi -lbcm2835

all: server

RasPi.o: $(RH)/RHutil/RasPi.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) $<

RH_RF69.o: $(RH)/RH_RF69.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) $<

RHSPIDriver.o: $(RH)/RHSPIDriver.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) $<

RHGenericDriver.o: $(RH)/RHGenericDriver.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) $<

RHHardwareSPI.o: $(RH)/RHHardwareSPI.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) $<

RHGenericSPI.o: $(RH)/RHGenericSPI.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) $<

server.o: main.cpp
	$(CC) $(CFLAGS) -c $(INCLUDE) -o server.o $<

server: RasPi.o RH_RF69.o RHSPIDriver.o RHGenericDriver.o RHHardwareSPI.o RHGenericSPI.o server.o
	$(CC) $^ $(LIBS) -g -o server

clean:
	rm -f *.o server

.PHONY: server.o server
