
## 欢迎来到 data-platform

### data-platform 是什么

data-platform 是一个**快速对接数**据源, 可以实现自由的加工(js代码)和整合数据， 并自动提供封装api的一个平台，并具有结构化的log输出和rest的接口风格；

其分为五个服务: 数据源对接服务(data-source)，数据源http代理服务(data-source-api), 数据api封装服务(data-api), 数据变量加工服务(data-derive),  数据平台总代理(data-proxy)；

 在使用平台的时候可以根据需求选择需要的服务, 不需要所有服务都使用;

数据接入和编写加工逻辑等, 可以通过 **data-platform-admin** 项目进行页面可视化编辑;

### 数据源对接服务

数据源对接服务, 可以快速的实现数据对接, 并提供rpc接口供访问

支持数据种类:  mysql， oracle，redis, couchdb， tidb，mongo，file，http

### 数据源http代理服务

对数据源对接服务得rpc接口进行封装，进行反向代理，提供http服务供调用，用户可以根据需求选择是否使用

### 数据api封装服务

对数据源接入的数据返回进行加工，返回加工后的结果，并自动提供http接口供访问，支持对访问参数做加工，对返回结果做加工， 对log输出结果做加工；

### 数据变量加工服务

对批量的数据api封装的结果进行整合，并加工和封装，并可以以变量和变量集合的方式作为接口调用的返回，支持对访问参数做加工，对返回结果做加工， 对log输出结果做加工，对变量整合成集合等；

### 数据平台总代理

对调用进行代理，规整化请求方法，清洗接口调用方式

## 安装方法(提前安装docker，docker-compose)

1. 拉取代码

   ```shell
   git clone git@git.100credit.cn:blockchain/data-platform.git
   ```

2. 制作二进制程序(如果已经存在可以跳过)

   ```shell
   cd  docker-deplog
   go build -o data-source ../services/main-control/main_data_source_plus.go   # 编译data-source
   go build -o data-source-api ../apis/data-source/api.go   # 编译data-source-api
   go build -o data-api ../webs/main-control/main_data_api.go # 编译data-api
   go build -o data-derive ../webs/main-control/main_data_derive.go # 编译data-derive
   go build -o data-proxy ../data-proxy/main.go   # 编译data-proxy
   ```

3. 修改配置文件

   配置文件位置： docker-compose/config 目录 ，  可以根据需求自己更改配置文件的内容

4. 启动数据平台依赖

   ```shell
   bash import
   bash deploy.sh start db
   ```

5. 启动数据平台服务

   ```shell
   bash deploy.sh start service
   ```

## 使用方法

1. 数据源接入，加工数据等方式

   数据源接入，加工数据，参数，log日志等都是通过页面配置的方式进行，改部分参考 **data-platform** 项目

1. 数据调用

   数据调用都是通过http的方式进行，暴露的接口集合如下

   #### 数据源相关接口

   ```
   数据源所需参数获取方式
   url: /ds/query/params
   method: POST
   params: code 数据源代码 
   
   数据调用方式
   url: /ds/query/data
   method: POST
   params: code 数据源代码 | data 数据源配置所需要的参数  | cache 0/1 是否需要缓存数据
   ```

   #### 数据api封装相关接口

   ```
   api所需参数获取方式
   url: /data-api/data/params
   method: POST
   params: code 数据api代码
   
   api数据调用方式
   url: /data-api/data/query
   method: POST
   params: code api代码 | data  api所需要的参数  | cache 0/1 是否需要缓存数据
   
   api接口验证方式
   url: /data-api/data/verify
   method: POST
   params: code api代码 
   ```

   #### 数据变量相关接口

   ```
      变量所需参数获取方式
      url: /data-derive/data/params
      method: POST
      params: code 数据变量代码
      
      变量调用方式
      url: /data-derive/data/queryDerive
      method: POST
      params: code 变量代码 | data 变量配置所需要的参数  | cache 0/1 是否需要缓存数据
      
      变量数据接口验证方式
      url: /data-api/data/verifyDerive
      method: POST
      params: code  变量代码 
      
      变量集调用方式
      url: /data-derive/data/queryDeriveSet
      method: POST
      params: code 变量集代码 | data 变量集配置所需要的参数  | cache 0/1 是否需要缓存数据
      
      变量集数据接口验证方式
      url: /data-api/data/verifyDeriveSet
      method: POST
      params: code  变量集代码 
   ```

   #### consul 服务取消

   ```
   curl --request PUT http://127.0.0.1:8500/v1/agent/service/deregister/{}
   ```

# 贡献代码

1. 拉取代码

2. 创建非maser分支，更改代码后提交请求

