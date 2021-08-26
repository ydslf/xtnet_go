# xtnet_go

#### 介绍


#### 软件架构
Tcp包结构
| pktHead |                      pktBody                     |    包结构
| pktLen  | msgID |              msgBody                     |    客户端消息(网关做router逻辑,自己判断消息是发送给哪个前端服务器; 或者没有网关)
| pktLen  | msgDirection | msgID |          msgBody          |    客户端消息(网关不做router逻辑,根据msgDirection判断发送给哪个前端服务器)
| pktLen  | rpcType | contextID  | msgType | msgID | msgBody |    服务器内部消息

websocket包结构
相对于Tcp包结构，没有pkgLen


服务器构架
           manager                              server manager  
     center     matching                        backend server

login    lobby  lobby     game  game            frontend server

    gate    gate    gate                        gate server

        client  client                          client


1. 通用序列号工具
2. netAgent 注册消息函数，能自动反序列化初函数的参数
3. rpc完善
4. 网络库完善
5. 加一种计时器，并且加上crontab

#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
