import os
from time import time_ns
import logging
import json
from smal.bot import SMALBot

# setup logging, we want timestamps
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s.%(msecs)03d %(levelname)s %(name)s - %(funcName)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger(__name__)
logger.setLevel(level=logging.INFO)


DEFAULT_PREFIX = "!"


class DemoBot(SMALBot):

    def on_event(self, evt):
        print("Got an event: ", e)

    def on_message(self, msg):

        if msg["type"] != "m.room.message":
            # not a room message
            logger.error(f"not a room message: {msg}")
            return

        if msg["sender"] == "get own id from missing code":
            # ignore own messages
            # for now just do not send valid commands by yourself
            logger.info(f"ignore own message: {msg}")
            return

        if "msgtype" in msg["content"].keys() and msg["content"]["msgtype"] != "m.text":
            # only react to messages, not emotes
            logger.debug(f"ignore unknown message type: {msg}")
            return

        if msg["content"]["body"] == "!stop":
            logger.info(f"stopping the bot")
            self.stop()
            return

        if msg["content"]["body"] == "!leave":
            logger.info(f"leaving room {msg['roomid']}")
            self.leaveroom(msg["roomid"])
            return

        if msg["content"]["body"].startswith("!echo"):

            txt = msg["content"]["body"][5:].strip()

            if txt == "":
                txt = "Empty text? Are you kidding me?"

            if msg["is_direct"]:
                self.sendmessage(msg["roomid"], txt)
            else:
                self.sendmessagereply(msg["roomid"], msg["id"], msg["sender"], txt)
            return

        logger.info(f"ignored a message: {msg}")


def main():
    # create and run the bot
    bot = DemoBot(DEFAULT_PREFIX)
    bot.run()


if __name__ == "__main__":
    main()
