#!/usr/bin/env python3

import adafruit_rfm9x
import argparse
import asyncio
import board
import busio
import digitalio
import json
import logging
import requests
import websockets

from enum import Enum
from typing import Any, Dict, List, Tuple

log = logging.getLogger(__name__)


def get_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()

    parser.add_argument('--server', '-s', required=True, type=str)
    parser.add_argument('--bridge-id', '-i', default=0, type=int)

    return parser.parse_args()


class RadioCommand(Enum):
    NOP = 0
    DISCOVER = 1
    STATUSUPDATE = 2


class ModeCommand(Enum):
    NOP = 0
    CHANGECOLOR = 1
    VOTEENABLE = 2
    VOTEDISABLE = 3

    @classmethod
    def translate_vote_mode(cls, value: int) -> 'ModeCommand':
        if value == 1:
            return cls.VOTEENABLE
        else:
            return cls.VOTEDISABLE


class Bridge:
    def __init__(self):
        self.args = get_args()
        self._loop = asyncio.get_event_loop()

        self._outgoing_message_queue = asyncio.Queue()

        self._radio = None

    @staticmethod
    def _rgb_to_bytestream(rgb: int) -> List[int]:
        r = (rgb & 0xFF0000) >> 16
        b = (rgb & 0x00FF00) >> 8
        g = (rgb & 0x0000FF)

        return [r, g, b]

    @staticmethod
    def _make_packet(device: Dict[str, Any]) -> Tuple[bytes, bytes]:
        color_pkt = []
        vote_pkg = []

        try:
            # Create color sequence packet.
            color_pkt.extend([
                RadioCommand.STATUSUPDATE.value,
                ModeCommand.CHANGECOLOR.value
            ])

            for rgb in device['colorsequence']:
                color_pkt.extend(self._rgb_to_bytestream(rgb))

            # Create voting packet
            vote_pkt.extend([
                RadioCommand.STATUSUPDATE.value,
                ModeCommand.translate_vote_mode(device['mode']).value
            ])

            return bytes(color_pkt), bytes(vote_pkt)
        except KeyError as e:
            log.error(e)
        except ValueError as e:
            log.error(e)
        finally:
            return bytes(), bytes()

    async def ws_listener(self) -> None:
        async with websockets.connect(
            f'ws://{args.server}/api/ws/bridge', ping_interval=None
        ) as ws:
            while True:
                payload = await ws.recv()
                content = json.loads(payload)

                log.debug(f'Got payload: {content}')

                for device in content:
                    color_pkt, vote_pkt = self._make_packet(device)

                    try:
                        if color_pkt:
                            await self._outgoing_message_queue.put({
                                'id': device['id'],
                                'packet': color_pkt
                            })

                        if vote_pkt:
                            await self._outgoing_message_queue.put({
                                'id': device['id'],
                                'packet': vote_pkt
                            })
                    except KeyError as e:
                        log.error(e)

    async def autodiscover(self) -> None:
        while True:
            await self._outgoing_message_queue.put({
                'packet': bytes([RadioCommand.DISCOVER.value, ModeCommand.NOP.value])
            })

            await asyncio.sleep(30)

    async def packet_send(self) -> None:
        if not self._radio:
            return

        while True:
            message = await self._outgoing_message_queue.get()
            if 'id' in message:
                self._radio.send(
                    message['packet'],
                    tx_header=(message['id', self.args.bridge_id, 0, 0])
                )
            else:
                self._radio.send(message['packet'])

    async def packet_recieve(self) -> None:
        while True:
            payload = self._radio.receive(rx_filter=self.args.bridge_id)

            if payload[0] == RadioCommand.DISCOVER.value:
                res = requests.get(f'http://{self.args.server}/api/device/all')
                content = json.loads(rest.text)

                top_id = 0
                for device in content:
                    if device['id'] > top_id:
                        top_id = device['id']

                requests.post(
                    f'http://{self.args.server}/api/device',
                    data={'id': top_id + 1, 'mode': 0, 'color': 0}
                )
            elif payload[0] == RadioCommand.STATUSUPDATE.value:
                res = requests.get(f'http://{self.args.server}/api/votes')
                content = json.loads(res.text)

                if payload[2] == 1:
                    content['positive'] += 1
                else:
                    content['negative'] += 1

                requests.post(f'http://{self.args.server}/api/votes', data=content)

    def run(self) -> None:
        # Setup radio
        cs = digitalio.DigitalInOut(board.CE1)
        rst = digitalio.DigitalInOut(board.D25)
        spi = busio.SPI(board.SCK, MOSI=board.MOSI, MISO=board.MISO)
        self._radio = adafruit_rfm9x.RFM9x(spi, cs, rst, 915.0)

        log.info('Radio initialized')

        asyncio.create_task(ws_listener())
        asyncio.create_task(autodiscover())
        asyncio.create_task(packet_send())
        asyncio.create_task(packet_receive())

        self._loop.run_forever()


if __name__ == '__main__':
    bridge = Bridge()
    bridge.run()
