# xtnet_go

#### 介绍


#### 软件架构
Tcp包结构
| pkgHead |                      pkgBody                |    包结构
| pkgLen  | msgID |              msgBody                |    客户端消息(网关做router逻辑,自己判断消息是发送给哪个前端服务器; 或者没有网关)
| pkgLen  | msgDirection | msgID |          msgBody     |    客户端消息(网关不做router逻辑,根据msgDirection判断发送给哪个前端服务器)
| pkgLen  | rpcType | contextID  | msgID |  msgBody     |    服务器内部消息(直接发送)

websocket包机构
相对于Tcp包结构，没有pkgLen


服务器构架
           manager                                 
     center     matching                        backend server

login    lobby  lobby     game  game            frontend server

    gate    gate    gate                        gate server

        client  client                          client

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
