import yaml
import asyncio
from nats.aio.client import Client as NATS
from nats.aio.errors import ErrTimeout, ErrNoServers

with open("./etc/default.yml", 'r') as stream:
    cfg = yaml.safe_load(stream)['nats']


with open("./desciption/model.json", 'r') as stream:
    model_data = stream.read()

async def publish_message():
    nc = NATS()

    # Connect to the NATS server
    await nc.connect(f"{cfg['address']}", max_reconnect_attempts=2)

    # Define the JetStream subject and message data
    subject = cfg['stream'].replace('*', 'test')
    message_data = model_data.encode()

    try:
        # Publish the message to JetStream
        await nc.publish(subject, payload=message_data)
        print(f"Published message: {message_data.decode()} to subject: {subject}")
    except ErrTimeout:
        print("Publish request timed out")
    except ErrNoServers:
        print("No NATS servers available")
    except Exception as e:
        print(e)

    # Close the connection
    await nc.close()

if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(publish_message())
