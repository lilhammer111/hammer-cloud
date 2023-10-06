# hammer-cloud
## 项目介绍
go web项目，使用httprouter和gin框架实现了两版，用了阿里云oss对象存储服务
## 项目演示
暂无
## 组织架构
```
.
├── cache  // redis缓存
├── common // 常量
├── config // OSS和rabbitmq相关配置
├── db // mysql数据库连接与接口
├── doc // sql 建表语句
├── go.mod 
├── go.sum
├── handler // 业务逻辑处理函数
├── meta // 文件源数据
├── middleware // 鉴权中间件
├── mq // rabbitmq
├── README.md
├── service // 路由及主入口函数
├── static // 静态文件
├── store // 对象存储接口
├── test 
└── util // 加密及通用响应
```

## 技术选型
|技术|说明|官网|
|-|-|-|
|gin|go web应用开发框架|https://gin-gonic.com/|
|mysql|关系型数据库|https://dev.mysql.com/doc/refman/8.0/en/|
|redis|内存数据存储|https://redis.uptrace.dev/guide/|
|OSS|对象存储服务|https://www.alibabacloud.com/zh/product/object-storage-service|

## 架构图
紧急生产中。。。

## 环境搭建
紧急搭建中。。。

## 许可证
紧急许可中。。。
