# TcpProxy
一个支持主备/哈希/随机分配模式的负载均衡代理服务器

负载均衡策略: 
    第一级采用来源IP计算hash值来分配到对应的组;
    第二级在组内采用配置的方式(master/hash/rand)分配到对应的节点

## 特性
1.　分配策略丰富，支持二级、多种分配策略

2.　性能强悍（见后续性能测试）

3.　自动检查后端服务器状态


## 获取
```
go get -u github.com/yl365/TcpProxy
```

## 使用
1. 修改配置文件config.json
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
2. 在配置目录运行: nohup ./TcpProxy &

## 性能测试
后端服务器采用spark作为目标服务器: `./spark -port 52241 "<h1>Ooops</h1>" `

1. 直连spark测试:
```
[yl@mobile-server-61 ~]$ ab -k -c 1000 -n 1000000 http://10.15.107.61:52241/
This is ApacheBench, Version 2.0.40-dev <$Revision: 1.146 $> apache-2.0
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Copyright 2006 The Apache Software Foundation, http://www.apache.org/

Benchmarking 10.15.107.61 (be patient)
Completed 100000 requests
Completed 200000 requests
Completed 300000 requests
Completed 400000 requests
Completed 500000 requests
Completed 600000 requests
Completed 700000 requests
Completed 800000 requests
Completed 900000 requests
Finished 1000000 requests


Server Software:        
Server Hostname:        10.15.107.61
Server Port:            -13295

Document Path:          /
Document Length:        14 bytes

Concurrency Level:      1000
Time taken for tests:   21.298565 seconds
Complete requests:      1000000
Failed requests:        0
Write errors:           0
Keep-Alive requests:    1000000
Total transferred:      154121044 bytes
HTML transferred:       14011004 bytes
Requests per second:    46951.52 [#/sec] (mean)
Time per request:       21.299 [ms] (mean)
Time per request:       0.021 [ms] (mean, across all concurrent requests)
Transfer rate:          7066.58 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   2.2      0     124
Processing:     0   20  13.1     19     119
Waiting:        0   20  13.1     19     119
Total:          0   20  13.3     19     148

Percentage of the requests served within a certain time (ms)
  50%     19
  66%     24
  75%     28
  80%     30
  90%     38
  95%     45
  98%     53
  99%     60
 100%    148 (longest request)
 ```
 2. 通过TcpProxy测试:
 ```
[yl@mobile-server-61 ~]$ ab -k -c 1000 -n 1000000 http://10.15.107.61:18080/
This is ApacheBench, Version 2.0.40-dev <$Revision: 1.146 $> apache-2.0
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Copyright 2006 The Apache Software Foundation, http://www.apache.org/

Benchmarking 10.15.107.61 (be patient)
Completed 100000 requests
Completed 200000 requests
Completed 300000 requests
Completed 400000 requests
Completed 500000 requests
Completed 600000 requests
Completed 700000 requests
Completed 800000 requests
Completed 900000 requests
Finished 1000000 requests


Server Software:        
Server Hostname:        10.15.107.61
Server Port:            18080

Document Path:          /
Document Length:        14 bytes

Concurrency Level:      1000
Time taken for tests:   20.702664 seconds
Complete requests:      1000000
Failed requests:        0
Write errors:           0
Keep-Alive requests:    1000000
Total transferred:      154000000 bytes
HTML transferred:       14000000 bytes
Requests per second:    48302.96 [#/sec] (mean)
Time per request:       20.703 [ms] (mean)
Time per request:       0.021 [ms] (mean, across all concurrent requests)
Transfer rate:          7264.28 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   3.2      0     171
Processing:     0   20   7.0     19      90
Waiting:        0   20   7.0     19      89
Total:          0   20   7.8     19     218

Percentage of the requests served within a certain time (ms)
  50%     19
  66%     22
  75%     24
  80%     25
  90%     29
  95%     32
  98%     36
  99%     40
 100%    218 (longest request)
```

欢迎试用并提出意见建议. 如果发现bug, 请Issues, 谢谢!
