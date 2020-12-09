# bpfly - BPF UI


## Requirements
- bcc-tools
- python-bcc
- linux-headersI

## Build

### Backend
```go build .```

### Frontend
```make -C elm```

## Run

### Backend

- Run a redis server:
```redis-server```

- And [webis](https://github.com/nicolasff/webdis)

- BPFly
```sudo ./bpfly```

- Open in Browser: elm/index.html
