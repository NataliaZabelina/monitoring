FROM ubuntu:18.04

WORKDIR /monitoring

COPY ./bin/monitoring ./bin/
COPY ./configs/config.json ./configs/

RUN apt-get update \
    && apt-get -y install sysstat \
    && apt-get -y install tzdata \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 50051
ENTRYPOINT ["./bin/monitoring"]
CMD ["grpc_server", "-c", "./configs/config.jso"]