import logging
import json
import time
import requests

# data={}
ts = time.time()
for i in range(1_000):
    res = requests.put("http://127.0.0.1:3000/put", data=json.dumps({str(i): str(i)}))
te = time.time()
print(f"Writing took {te-ts}s")

payload_data = {
    "key": "foo"
}

json_data = json.dumps(payload_data)

headers = {
    "Content-Type": "application/json"
}

ts = time.time()
response = requests.get("http://127.0.0.1:3000/get", data=json.dumps(payload_data), headers={"Content-Type": "application/json"})
te = time.time()
print(f"{response.text} exec time: {te-ts}s")