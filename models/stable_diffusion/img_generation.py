import base64
import os
import time
from diffusers import StableDiffusionPipeline
import torch

model_id = "runwayml/stable-diffusion-v1-5"
pipe = StableDiffusionPipeline.from_pretrained(model_id, torch_dtype=torch.float16)
pipe = pipe.to("cuda")

def generate_img(prompt: str):
  image = pipe(prompt).images[0]
  timestamp = time.time()
  file_name = f"generated_{timestamp}.png"
  path = "../data/"
  image.save(path + file_name)
  print(f"Path is: {path + file_name}")

  with open(path + file_name, "rb") as image_file:
    encoded_string = base64.b64encode(image_file.read())

  base64_img = encoded_string.decode('utf-8')
  os.remove(path=path + file_name)
  return (file_name, base64_img)
  
