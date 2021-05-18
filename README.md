# System monitoring

Демон - программа, собирающая информацию о системе, на которой запущена, и отправляющая её своим клиентам по GRPC.

## Запуск

Сервер:
 
 - make run_server
 - make run_client

 - ./monitoring grpc_server --config=configs/config.json
 - ./monitoring grpc_client --server=":50051" --timeout=5 --period=15

 Подробнее: 
 - ./monitoring --help
 - ./monitoring grpc_client --help
 - ./monitoring grpc_server --help

 Для запуска сервера в контейне:
 - make run_img

 ## Пример конфигурации ceрвера

 ```
 {
    "host": "localhost",
    "port": 50053,
    "log": {
        "level": "debug",
        "file": "./m.log"
    },
    "collector": {
        "timeout": 1,
        "statistics": {
            "load_system": true,
            "load_cpu": true,
            "load_disk": true,
            "top_talkers": false,
            "stat_network": false
        }
    }
}
```
