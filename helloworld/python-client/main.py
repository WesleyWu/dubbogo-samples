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

from fastapi import FastAPI
from grpc import Channel, insecure_channel
from uvicorn import Server, Config

from proto.helloworld_pb2 import HelloRequest
from proto.helloworld_pb2_grpc import GreeterStub

channel: Channel
stub: GreeterStub
service_endpoint = ""
xds_bootstrap_json_path = "xds_bootstrap.json"


#
# # Import the message definitions
# from envoy.service.status.v3 import csds_pb2
# # Import the gRPC service and stub
# from envoy.service.status.v3 import csds_pb2_grpc


@asynccontextmanager
async def lifespan(app: FastAPI):
    startup()

    yield

    shutdown()


def generate_xds_bootstrap(pod_ip, pod_name, pod_namespace):
    # read input file
    fin = open(xds_bootstrap_json_path, "rt")
    # read file contents to string
    data = fin.read()
    # replace all occurrences of the required string
    data = data.replace('{POD_IP}', pod_ip)
    if pod_name is not None:
        data = data.replace('{POD_NAME}', pod_name)
    if pod_namespace is not None:
        data = data.replace('{POD_NAMESPACE}', pod_namespace)
    # close the input file
    fin.close()
    # open the input file in write mode
    fin = open(xds_bootstrap_json_path, "wt")
    # overwrite the input file with the resulting data
    fin.write(data)
    # close the file
    fin.close()


def startup():
    logging.basicConfig(encoding='utf-8', level=logging.INFO)
    global channel, stub, service_endpoint

    service_endpoint = os.environ.get('SERVICE_ENDPOINT')

    if service_endpoint is None:
        raise Exception("SERVICE_ENDPOINT not set")

    pod_ip = os.environ.get('POD_IP')
    pod_name = os.environ.get('POD_NAME')
    pod_namespace = os.environ.get('POD_NAMESPACE')
    logging.info(f'service_endpoint: {service_endpoint}, '
                 f'pod_ip: {pod_ip}, '
                 f'pod_name: {pod_name}, '
                 f'pod_namespace: {pod_namespace}')

    # check if k8s deployment
    if pod_ip is not None:
        # replace envoy xds bootstrap file content
        generate_xds_bootstrap(pod_ip, pod_name, pod_namespace)
        channel = insecure_channel(service_endpoint)
    else:
        # calling local endpoint
        channel = insecure_channel(service_endpoint)
    stub = GreeterStub(channel)
    logging.info(f'Opened rpc {service_endpoint} channel ...')


def shutdown():
    global channel, stub, service_endpoint
    if channel:
        logging.info(f"Closing rpc {service_endpoint} channel")
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
