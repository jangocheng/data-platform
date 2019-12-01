#! /bin/bash

platform_image_name="data-platform-env"
platform_admin_image_name="data-platform-admin"
micro_image_name="micro-client"
mysql_image_name="mysql"
redis_image_name="redis"
consul_image_name="consul"

images_slice=("$platform_image_name" "$micro_image_name" "$mysql_image_name" "$redis_image_name" \
                "$consul_image_name" "$platform_admin_image_name")


import_image(){
    echo "同步基础镜像中......"
    for image_name in ${images_slice[@]}
    do
        image_exist=$(docker images | grep "$image_name" | grep "latest")
        if [ ! "$image_exist" ]; then
            docker import "./images-tar/$image_name.tar" $image_name
            if [ $? -ne 0 ]; then
                return 1
            fi
        fi
    done
    echo "镜像同步完成......"
    return 0
}

import_image
if [ $? -ne 0 ]; then
    echo "基础镜像同步出现异常"
    exit 1
fi

