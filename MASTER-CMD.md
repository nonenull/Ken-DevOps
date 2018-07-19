MASTER-CMD
======

## 简介
    master的命令行客户端

## 目录结构
    ken-master      Master端代码
       | -- bin     保存编译脚本, 编译文件等
       | -- certs   保存Servant 公钥目录
       | -- conf    配置目录
       | -- logs    日志目录
       | -- src     源码目录

## 使用说明

### Usage
    Usage: ./command [-S] [servantID] [func] [args]
    
    Example:
        ./command -S nginxserver network.getip -i eth0
    Option:
        [-S]	(可选)指定处理的连接类型为短连接, 默认为长连接
        [servantID]	servant的主机名
        [func]	在servant主机上执行的函数
        [args]	传递给执行函数的参数