#!/usr/bin/env python
# Test MQTT Task Queue
import os
import glob
import json

# check if pika and dotenv are installed
try:
    import pika
except ImportError:
    print("pika not installed")
    print("use 'pip install pika' to install")
    exit(1)
try:
    from dotenv import load_dotenv
except ImportError:
    print("dotenv not installed")
    print("use 'pip install python-dotenv' to install")
    exit(1)

def must_load_amqp_url():
    # find and load all env from parent folder
    parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    env_files = glob.glob(os.path.join(parent_dir, "..", "**", ".env.*"), recursive=True)
    for env_file in env_files:
        # print("Loading env file:", env_file)
        load_dotenv(env_file)

    # assert AMQP_URL is set
    AMQP_URL = os.getenv("AMQP_URL")
    assert AMQP_URL, "AMQP_URL must be set"
    return AMQP_URL

# connect to AMQP
AMQP_URL = must_load_amqp_url()
# connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
connection = pika.BlockingConnection(pika.URLParameters(AMQP_URL))
channel = connection.channel()

# tasks queue
channel.queue_declare(queue="ci_tasks", durable=True)

def send(message):
    channel.basic_publish(
        exchange="",
        routing_key="ci_tasks",
        body=message,
        properties=pika.BasicProperties(delivery_mode=pika.DeliveryMode.Persistent),
    )
    print(f" [x] Sent {message}")

# send a message
# send("Hello World!")
send(json.dumps({"git": "https://github.com/CyperpunksAmurai/go-build-it-test-project"}))
# close connection
connection.close()
