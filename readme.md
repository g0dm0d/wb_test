# Performance

### Nginx


```
Running 30s test @ http://127.0.0.1:8080/api/b563feb7b2b84b6test
  8 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    36.17ms   39.41ms 299.98ms   87.22%
    Req/Sec     1.89k     1.33k    4.16k    50.25%
  453284 requests in 30.10s, 437.90MB read
Requests/sec:  15061.08
Transfer/sec:     14.55MB
```

### Directly

```
Running 30s test @ http://localhost:3000/b563feb7b2b84b6test
  8 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.29ms    2.54ms  58.85ms   78.33%
    Req/Sec    16.32k     2.15k   27.85k    69.18%
  3903431 requests in 30.10s, 3.52GB read
Requests/sec: 129682.30
Transfer/sec:    119.59MB
```

# Installation

```shell
cp ./configs/.env.example .env
source .env
docker compose build
docker compose up
```

in browser http://localhost:8080 simple web
