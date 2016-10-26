# TcpProxy
一个支持主备/哈希/随机分配模式的负载均衡代理服务器

【TcpProxy技术交流QQ群】: 99328252　　加入群：http://jq.qq.com/?_wv=1027&k=40qDFxw

```
负载均衡策略: 
    第一级采用来源IP计算hash值来分配到对应的组;
    第二级在组内采用配置的方式(master/hash/rand)分配到对应的节点
```

## 特性
1.　分配策略丰富，支持二级、多种分配策略

2.　性能强悍（见后续性能测试）

3.　自动检查后端服务器状态


## 获取
```
go get -u github.com/yl365/TcpProxy
```

## 使用

1、 修改配置文件config.json

```
{
    "Listen": ":18080",             //监听端口
    "Mode": "master/hash/rand",     //分配模式,master/hash/rand根据场景选择一种
    "AllHost": [
        {
            "min": 0,               //每组服务器的处理范围[0--999],每组分配不要重合
            "max": 333, 
            "Hosts": [
                {
                    "IP": "127.0.0.1:2000", 
                    "status": 0
                }, 
                {
                    "IP": "127.0.0.1:2001", 
                    "status": 0
                }
            ]
        }, 
        {
            "min": 334, 
            "max": 666, 
            "Hosts": [
                {
                    "IP": "127.0.0.1:2002", 
                    "status": 0
                }, 
                {
                    "IP": "127.0.0.1:2003", 
                    "status": 0
                }
            ]
        }, 
        {
            "min": 667, 
            "max": 999, 
            "Hosts": [
                {
                    "IP": "127.0.0.1:2004", 
                    "status": 0
                }, 
                {
                    "IP": "127.0.0.1:2005", 
                    "status": 0
                }
            ]
        }
    ]
}
```

2、 在配置目录运行: nohup ./TcpProxy &


## 性能测试

后端服务器采用spark作为目标服务器: `./spark -port 52241 "<h1>Ooops</h1>" `

1、 直连spark测试:

``` 
[yl@mobile-server-61 bin]$ ./gobench -c 3000 -k -t 10  -u http://10.15.107.61:52241/
Dispatching 3000 clients
Waiting for results...

Requests:                           729340 hits
Successful requests:                729339 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:            72933 hits/sec
Read throughput:                   9481420 bytes/sec
Write throughput:                  6587163 bytes/sec
Test time:                              10 sec
```

2、 通过haproxy测试:

```
[yl@mobile-server-61 bin]$ ./gobench -c 3000 -k -t 10  -u http://10.15.107.61:8008/
Dispatching 3000 clients
Waiting for results...

Requests:                           128291 hits
Successful requests:                128291 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:            12829 hits/sec
Read throughput:                   1667952 bytes/sec
Write throughput:                  1155709 bytes/sec
Test time:                              10 sec
```

3、 通过TcpProxy测试:

```
[yl@mobile-server-61 bin]$ ./gobench -c 3000 -k -t 10  -u http://10.15.107.61:18080/
Dispatching 3000 clients
Waiting for results...

Requests:                           403622 hits
Successful requests:                403622 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:            40362 hits/sec
Read throughput:                   5247086 bytes/sec
Write throughput:                  3659598 bytes/sec
Test time:                              10 sec
```

欢迎试用并提出意见建议。如果发现bug，请Issues，谢谢！

如果你觉得这个项目有意义，或者对你有帮助，或者仅仅是为了给我一点鼓励，请不要吝惜，给我一个*star*，谢谢！≡ω≡

【TcpProxy技术交流QQ群】: 99328252　　加入群：http://jq.qq.com/?_wv=1027&k=40qDFxw
