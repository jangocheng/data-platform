#! /bin/bash

cmd=$1
pro=$2

services_start_db(){
    echo "容器启动中......"
    docker-compose -f docker-compose-db.yml up -d
    if [ $? -ne 0 ]; then
        return 1
    fi
    echo "容器启动完成....."
}

services_start_service(){
    echo "容器启动中......"
    docker-compose -f docker-compose-local.yml up -d
    if [ $? -ne 0 ]; then
        return 1
    fi
    echo "容器启动完成....."
}

services_stop(){
    echo "容器停止中......"
        docker-compose -f docker-compose-local.yml down
        if [ $? -ne 0 ]; then
            return 1
        fi
    echo "容器停止完成....."
}
db_stop(){
    echo "容器停止中......"
        docker-compose -f docker-compose-db.yml down
        if [ $? -ne 0 ]; then
            return 1
        fi
    echo "容器停止完成....."
}


if [ "$cmd" == "start" ]; then
    if [ "$pro" == "db" ]; then
        services_start_db
        if [ $? -ne 0 ]; then
            echo "服务启动异常"
            exit 1
        fi
    fi
    if [ "$pro" == "service" ]; then
        services_start_service
        if [ $? -ne 0 ]; then
            echo "服务启动异常"
            exit 1
        fi
    fi
elif [ "$cmd" == "stop" ]; then
    if [ "$pro" == "service" ]; then
        services_stop
        if [ $? -ne 0 ]; then
            echo "服务停止异常"
            exit 1
        fi
    fi
    if [ "$pro" == "db" ]; then
        db_stop
        if [ $? -ne 0 ]; then
            echo "服务停止异常"
            exit 1
        fi
    fi
fi 
