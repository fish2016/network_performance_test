# 介绍
此项目为测试http网络程序的性能，初始代码克隆自项目[zap](https://github.com/zigzap/zap)的wrk目录


原项目，支持对zip、go、rust等网络框架性能测试，在此基础上，计划加上对fasthttp、gnet等网络框架的性能测试


# 开始


## 准备工作

### 安装wrk

### 配置脚本权限
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


## 输出统计图表
``````
./wrk/measure_all.sh //需要先执行一遍全量性能测试（会产生测试xxx.perlog文件，graph.py会根据perlog文件输出图表）

python3 ./wrk/graph.py
``````

# 测试报告

## 设备1
### 硬件信息
========== CPU Info ==========

Model name:                           Intel(R) Core(TM) i5-10200H CPU @ 2.40GHz
Thread(s) per core:                   2
Core(s) per socket:                   4
Socket(s):                            1

========== Memory Info ==========

23GB RAM

### 性能报告
### 性能报告
![req_per_sec](https://github.com/fish2016/network_performance_test/blob/main/doc/test_result/req_per_sec_graph.png?raw=true)
![xfer_per_sec](https://github.com/fish2016/network_performance_test/blob/main/doc/test_result/xfer_per_sec_graph.png?raw=true)