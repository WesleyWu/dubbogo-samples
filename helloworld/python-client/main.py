# Copyright 2015 gRPC authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""The Python implementation of the GRPC helloworld.Greeter client."""

import logging
import os
from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI
from grpc import Channel, insecure_channel
from uvicorn import Server, Config

from proto.helloworld_pb2 import HelloRequest
from proto.helloworld_pb2_grpc import GreeterStub

channel: Channel
stub: GreeterStub
greeter_service_endpoint = ""


@asynccontextmanager
async def lifespan(app: FastAPI):
    startup()

    yield

    shutdown()


def startup():
    global channel, stub, greeter_service_endpoint
    try:
        greeter_service_endpoint = os.environ['GREETER_SERVICE_ENDPOINT']
    except KeyError:
        load_dotenv()
        greeter_service_endpoint = os.environ['GREETER_SERVICE_ENDPOINT']
    channel = insecure_channel(greeter_service_endpoint)
    stub = GreeterStub(channel)
    logging.basicConfig(encoding='utf-8', level=logging.INFO)
    logging.info(f'Opened rpc {greeter_service_endpoint} channel ...')


def shutdown():
    global channel, stub, greeter_service_endpoint
    if channel:
        logging.info(f"Closing rpc {greeter_service_endpoint} channel")
        channel.close()


app = FastAPI(lifespan=lifespan)


@app.get("/hello/{name}")
def say_hello(name):
    response = stub.SayHello(HelloRequest(name=name))
    logging.info(f"Got rpc response: {response.message}")
    return response.message


if __name__ == '__main__':
    # for development usage only.
    # remove reload=True for production or use helm k8s deployment
    s = Server(Config(host="0.0.0.0", port=5000, reload=True, app=app))
    s.run()
