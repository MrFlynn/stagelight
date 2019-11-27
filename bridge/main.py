import adafruit_rfm9x
import board
import busio
import digitalio
import time

console_message = """
Please select from one of the following modes:
0: Display color red.
1: Display color green.
2: Display color blue.
3: Go into voting mode.
"""


def main() -> None:
    chip_select = digitalio.DigitalInOut(board.CE1)
    reset = digitalio.DigitalInOut(board.D25)

    spi = busio.SPI(board.SCK, MOSI=board.MOSI, MISO=board.MISO)
    radio = adafruit_rfm9x.RFM9x(spi, chip_select, reset, 915.0)
    radio.tx_power = 23

    while True:
        try:
            print(console_message)
            mode = int(input("Please select a number [0-3]: "))

            radio.send(bytes([mode]), tx_header=(1, 0, 0, 0))
            print('Transmitting...')

            packet = radio.receive(with_header=True)

            if packet:
                print(f'Got vote: {packet}')
        except KeyboardInterrupt:
            print('\nGoodbye...')
            raise SystemExit(0)

if __name__ == '__main__':
    main()
