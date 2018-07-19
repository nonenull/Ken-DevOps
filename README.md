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
              | -- ken-common   主要保存一些公共包, 具体看包名便知
              | -- ken-master   Master端代码
                    | -- bin    保存编译脚本, 编译文件等
                    | -- certs  保存Servant 公钥目录
                    | -- conf   配置目录
                    | -- logs   日志目录
                    | -- src    源码目录
              | -- ken-servant  Servant端代码
                    | -- bin    保存编译脚本, 编译文件等
                    | -- cert   保存本机生成的公私密钥
                    | -- conf   配置目录
                    | -- logs   日志目录
                    | -- src    源码目录
              | -- ken-test     测试用代码

## 关键点说明

### Routers
    类似HTTP框架的使用方法, 需要为控制器注册路由
    
### Controller
    控制器, Servant端可调用的模块