# kdebug

> K8s debug pod — lightweight Go HTTP server for testing Kubernetes, APISIX, and network features.

[![CI](https://github.com/jory-callan/kdebug/actions/workflows/ci.yml/badge.svg)](https://github.com/jory-callan/kdebug/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/jory-callan/kdebug)](https://golang.org/dl/)
[![Docker Pulls](https://img.shields.io/badge/docker-ghcr.io/jory--callan/kdebug-blue?logo=github)](https://github.com/jory-callan/kdebug/pkgs/container/kdebug)
[![License](https://img.shields.io/github/license/jory-callan/kdebug)](LICENSE)

---

## Endpoints

| Route | Method | Description |
|-------|--------|-------------|
| `/ping` | GET | Health check → `{"code":0,"msg":"pong"}` |
| `/echo` | GET/POST | Echo query, body & headers |
| `/ip` | GET | Client real IP (respects X-Real-Ip / X-Forwarded-For) |
| `/env` | GET | Pod name, node name, version, start time |
| `/delay?ms=500` | GET | Simulate latency (default 100ms) |
| `/mem?mb=100&ms=10000` | GET | RAM spike test (default 1MB for 2s) |
| `/cpu?ms=5000&cores=2&percent=80` | GET | CPU burn test (default 2s, all cores, 80%) |
| `/` | GET | List all available routes |

## Quick Start

**Local:**
```bash
go run .                          # default: echo framework
go run . -c gin                   # gin framework
go run . -c mux                   # gorilla/mux
go run . -c http                  # net/http stdlib
```

**Docker:**
```bash
docker build -t kdebug .
docker run -p 8080:8080 kdebug
```

**K8s:**
```bash
kubectl apply -f k8s/kdebug.yml
```

## Features

- **4 frameworks in 1 binary**: echo (default), gin, gorilla/mux, net/http — switch with `-c`
- **Sub-10MB image**: multi-stage Docker build, Alpine-based
- **K8s native**: built-in liveness/readiness probes via `/ping`
- **Network testing**: client IP detection, header echo, latency simulation
- **Resource testing**: configurable CPU spikes and memory allocation
- **Multiple frameworks**: same endpoints, same JSON response shape — compare frameworks side by side

## Use Cases

- ☑ Test **K8s probes, rolling updates, and pod lifecycle**
- ☑ Validate **Ingress, APISIX routes, and traffic splitting**
- ☑ Simulate **latency, high CPU, and memory pressure** during load tests
- ☑ Verify **client IP preservation** through proxies and load balancers
- ☑ Replace production services for **canary/blue-green testing**

## Framework Comparison

```bash
# Start 4 instances on different ports:
go run . -c echo &                # :8080
go run . -c gin &                 # :8081  (flag.Parse uses first)
go run . -c mux &                 # :8082
go run . -c http &                # :8083
```

```bash
curl -s http://localhost:8080/ping  # echo
curl -s http://localhost:8081/ping  # gin
curl -s http://localhost:8082/ping  # mux
curl -s http://localhost:8083/ping  # stdlib
```

> All return `{"code":0,"msg":"pong"}`.

## Test Commands

```bash
# Liveness
curl http://localhost:8080/ping

# Echo
curl http://localhost:8080/echo -d 'hello=world'

# Client IP
curl -H 'X-Real-Ip: 1.2.3.4' http://localhost:8080/ip

# Latency: 500ms
curl http://localhost:8080/delay?ms=500

# Memory: allocate 200MB for 10s
curl http://localhost:8080/mem?mb=200&ms=10000

# CPU: 2 cores at 80% for 5s
curl "http://localhost:8080/cpu?ms=5000&cores=2&percent=80"
```

## License

MIT — use freely, no warranty.

---

> **中文简介**  
> kdebug 是一个超轻量的 Go HTTP 服务，专为 K8s 和 APISIX 功能验证设计。  
> 支持 Echo/Gin/Mux/NetHTTP 四种框架、客户端 IP 检测、延迟/内存/CPU 模拟。  
> Docker 镜像 < 10MB，开箱即用。
