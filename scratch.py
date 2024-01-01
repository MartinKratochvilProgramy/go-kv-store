import json
import requests

data = {"key": "1"}
res = requests.get("http://127.0.0.1:3000/", data=json.dumps(data))

print(res.text)