import onnxruntime as ort
import numpy as np
from PIL import Image
import huggingface_hub as hf
import io

from model_pb2_grpc import ImageGeneratorServicer
from model_pb2 import GenerateImageRequest, GenerateImageResponse

class ImageGenerator(ImageGeneratorServicer):
    def __init__(self, model_path: str, model_dir: str):
        self.download_model(model_id=model_path, model_dir=model_dir)
        self.session = ort.InferenceSession(model_dir + "/model.onnx", providers=['CPUExecutionProvider'])
        self.input_name = self.session.get_inputs()[0].name
        self.output_name = self.session.get_outputs()[0].name

    def GenerateImage(self, request, context):
        print(f"Received request: {request}")
        input_data = np.random.rand(1, 64).astype(np.float32)

        result = self.session.run([self.output_name], {self.input_name: input_data})[0]

        image_array = np.squeeze(result)
        image_array = ((image_array + 1) * 127.5 ).astype(np.uint8)
        img = Image.fromarray(image_array)
        print(f"img shape: {img.size}")
        img_byte_arr = io.BytesIO()
        img.save(img_byte_arr, format='PNG')
        img_byte_arr.seek(0)
        response = GenerateImageResponse(image_data=img_byte_arr.read())
        return response
    
    def download_model(self, model_id: str, model_dir: str):
        hf.hf_hub_download(repo_id=model_id, filename="model.onnx", repo_type="model", local_dir=model_dir)
