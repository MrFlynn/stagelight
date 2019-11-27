import adafruit_rfm9x
import asyncio
import board
import busio
import digitalio
import json
import sys
import time
import websockets


def main() -> None:
    chip_select = digitalio.DigitalInOut(board.CE1)
    reset = digitalio.DigitalInOut(board.D25)

    spi = busio.SPI(board.SCK, MOSI=board.MOSI, MISO=board.MISO)
    radio = adafruit_rfm9x.RFM9x(spi, chip_select, reset, 915.0)
    radio.tx_power = 23

    async def command_forwarder():
        async with websockets.connect(sys.argv[1]) as ws:
            while True:
                content = await ws.recv()
                devices = json.loads(content)

                for d in devices:
                    packet = None

                    color = d['colors'][0]
                    if color == 0xFF0000:
                        packet = 0
                    elif color == 0x00FF00:
                        packet = 1
                    elif color == 0x0000FF:
                        packet = 2

                    if packet is not None:
                        radio.send(bytes([packet]), tx_header=(1, 0, 0, 0))

    asyncio.get_event_loop().run_until_complete(command_forwarder())

if __name__ == '__main__':
    main()
