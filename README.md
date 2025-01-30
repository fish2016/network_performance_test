# 介绍
此项目为测试http网络程序的性能，初始代码克隆自项目[zap](https://github.com/zigzap/zap)的wrk目录


原项目，支持对zip、go、rust等网络框架性能测试，在此基础上，计划加上对fasthttp、gnet等网络框架的性能测试


# 开始


## 准备工作
``````
chmod 755 ./wrk/measure.sh
chmod 755 ./wrk/measure_all.sh
``````
## 单项测试

第二个参数为测试对象：

可选参数为： zig-zap go python python-sanic rust-bythebook rust-bythebook-improved rust-clean rust-axum csharp cpp-beast
``````
./wrk/measure.sh go
``````

## 全部测试
启动测试wrk目录下的所有测试对象

``````
./wrk/measure_all.sh
``````