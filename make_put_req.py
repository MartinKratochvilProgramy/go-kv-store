import time
import requests
import json
from concurrent.futures import ProcessPoolExecutor
from urllib3.util.retry import Retry
from requests.adapters import HTTPAdapter

session = requests.Session()
retry = Retry(connect=3, backoff_factor=0.5)
adapter = HTTPAdapter(max_retries=retry)
session.mount('http://', adapter)

def send_put_request(start, end):
    url = "http://127.0.0.1:3000/"
    for i in range(start, end):
        data = {str(i): str(i)}
        _ = session.put(url, data=json.dumps(data))

def send_get_request(start, end):
    url = "http://127.0.0.1:3000/"
    for i in range(start, end):
        data = {"key": str(i)}
        _ = session.get(url, data=json.dumps(data))

def main():
    N = 4_000
    processes = 8
    step = N // processes

    with ProcessPoolExecutor(max_workers=processes) as executor:
        ranges = [(i * step, (i + 1) * step) for i in range(processes)]

        put_futures = [executor.submit(send_put_request, start, end) for start, end in ranges]
        get_futures = [executor.submit(send_get_request, start, end) for start, end in ranges]

        ts = time.time()    
        for future in put_futures + get_futures:
            future.result()
        te = time.time()
        print(f"Exec time: {te-ts}s")

if __name__ == "__main__":
    main()
