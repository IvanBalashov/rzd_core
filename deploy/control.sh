#!/usr/bin/env bash

case $1 in
    "build" )
        docker build -t rzd_core:v0.1 ./
    ;;
    "run" )
        docker run --restart unless-stopped -d --name rzd_core_v0.1 rzd_core:v0.1
    ;;
    "stop" )
        docker stop rzd_core_v0.1
    ;;
    "rm" )
        docker rm rzd_core_v0.1
    ;;
    "rmi" )
        docker rmi rzd_core:v0.1
    ;;
    "rebuild" )
        docker stop rzd_core_v0.1
        docker rm rzd_core_v0.1
        docker rmi rzd_core:v0.1
        docker build -t rzd_core:v0.1 ./
        docker run --restart unless-stopped -d --name rzd_core_v0.1 rzd_core:v0.1
    ;;
    "inter")
        docker run -it --entrypoint /bin/bash --name rzd_core_v0.1 rzd_core:v0.1
    ;;
esac