from flask import Flask, request, jsonify
from paddleocr import PaddleOCR
import io
import json

app = Flask(__name__)
ocr = PaddleOCR(use_angle_cls=True, lang='ch')

@app.route('/ocr', methods=['POST'])
def ocr_endpoint():
    file = request.files['image']
    img_bytes = file.read()
    result = ocr.ocr(img_bytes, cls=True)
    print(json.dumps(result[0]))
    return jsonify(result[0])

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=9527)
