# compilesh-hello

A small HTTP server used to test compile.sh deployments end to end. It listens
on `$PORT` (set by the platform) and has endpoints that answer with chosen
status codes, so the console's Requests page has 2xx, 4xx and 5xx traffic.

| Endpoint | Answer |
| --- | --- |
| `GET /` | 200, the platform variables the app received |
| `GET /health` | 200 |
| `POST /items` | 201 with `{"name": "..."}`, 400 otherwise |
| `GET /bad-request` | 400 |
| `GET /error` | 500 |
| `GET /status/{code}` | any status from 200 to 599 |
| `GET /slow?ms=1500` | 200 after a delay (at most 10 s) |
| anything else | 404 |
