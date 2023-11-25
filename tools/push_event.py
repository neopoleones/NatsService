import yaml
import asyncio
import nats
from nats.js.api import StreamConfig, RetentionPolicy

with open("./etc/default.yml", 'r') as stream:
    cfg = yaml.safe_load(stream)['nats']


with open("./desciption/model.json", 'r') as stream:
    model_data = stream.read()


async def run():
    nc = await nats.connect("nats://127.0.0.1:4222")
    js = nc.jetstream()

    stream_config = StreamConfig(name=cfg['stream'].replace('.*', ''), subjects=[cfg['stream']])
    await js.add_stream(stream_config)
    await js.publish(cfg['stream'].replace('*', 'test'), model_data.encode())

    await nc.close()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run())
    loop.close()