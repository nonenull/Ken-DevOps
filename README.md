Ken-DevOps
======

## 简介
    一个 client -> server 的通信框架, 可用于运维自动化.
    可基于tls加密

## 目录结构
    Ken-DevOps
        |
        | -- src
              |
              | -- ken-common   主要保存一些公共包
              | -- ken-master   Master端代码
              | -- ken-servant  Servant端代码
              | -- ken-test     测试用代码

## 关键点说明

### Routers
    类似HTTP框架的使用方法, 需要为控制器注册路由
    
### Controller
    控制器, Servant端可调用的模块