```
Running 30s test @ http://172.20.0.4:3000/b563feb7b2b84b6test
  8 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     8.57ms   10.89ms 139.80ms   86.66%
    Req/Sec     9.32k     1.43k   17.37k    70.03%
  2229926 requests in 30.09s, 1.99GB read
Requests/sec:  74105.78
Transfer/sec:     67.78MB
```

# Installation

```shell
cp ./configs/.env.example .env
source .env
docker compose build
docker compose up
```

in browser http://localhost:8080 simple web
