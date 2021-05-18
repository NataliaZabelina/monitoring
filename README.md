# System monitoring

Демон - программа, собирающая информацию о системе, на которой запущена, и отправляющая её своим клиентам по GRPC.

## Запуск

 - ./monitoring grpc_server --config=configs/config.json
 - ./monitoring grpc_client --server=":50051" --timeout=5 --period=15

 ## Пример конфигурации

 ```
 {
    "host": "localhost",
    "port": 50053,
    "log": {
        "level": "debug",
        "file": "./mon.log"
    },
    "collector": {
        "timeout": 1,
        "statistics": {
            "load_system": true,
            "load_cpu": true,
            "load_disk": true,
            "top_talkers": true,
            "stat_network": false
        }
    }
}
```

