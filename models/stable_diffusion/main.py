import base64
import requests
from flask import Flask, request, jsonify
from img_generation import generate_img

app = Flask(__name__)

fileserver_create_img = "http://95.179.167.137:7676/api/v1/create/file"

@app.route("/")
def hello_world():
    return "<p>MilkyTeadrop Image Generation Api</p>"

@app.route('/api/v1/generate/', methods=['POST'])
def generate_img_api():
    if request.method == 'POST':
        json_data = request.json
        print(json_data)
        message = json_data.get('message')

        if not message:
            return handle_bad_request()

        print("Message to generate img for is: ", message)
        (file_path, base64_img) = generate_img(prompt=message)
        print("Name of file is: ", file_path)
        
        request_body = {
            "filename": file_path,
            "filetype": "image",
            "fileContent": base64_img
        }
        
        r = requests.post(url=fileserver_create_img, json=request_body)
        print("Response from fileserver is: ", r.status_code)
       
        return jsonify({
            "status": "ok",
            "message": "Image was generated",
            "filename": file_path
        }), 200
    else:
        return handle_bad_request()

def handle_bad_request():
    return jsonify({
        "status": "error",
        "message": "'message' is missing in request"
    }), 400

if __name__ == '__main__':
    app.run()