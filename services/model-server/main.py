import grpc
from concurrent import futures
import os

import model_pb2 as pb2
import model_pb2_grpc as pb2_grpc

from model import ImageGenerator

PORT = os.getenv("PORT")

MODEL_ID = os.getenv("MODEL_ID")

MODEL_DIR = os.getenv("MODEL_DIR", "/model")

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_ImageGeneratorServicer_to_server(
        ImageGenerator(model_path=MODEL_ID, model_dir=MODEL_DIR), server
    )
    server.add_insecure_port(f"[::]:{PORT}")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    print(f"Starting server on port {PORT} with model {MODEL_ID} in directory {MODEL_DIR}")
    serve()