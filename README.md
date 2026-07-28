# XTNet — Go 游戏服务器网络库

高性能、可扩展的 Go 语言游戏服务器网络层框架，支持 TCP/WebSocket 协议，内置路由与 RPC 机制。

---

## 简介

XTNet 是一个面向游戏服务器的网络通信库，提供：

- **TCP** 与 **WebSocket** 双协议支持
- 灵活的消息包结构设计
- 网关路由（Router）与无网关直连两种通信模式
- 服务端内部 RPC 通信
- 分层服务器架构（Gate / Frontend / Backend）

---

## 软件架构

### TCP 包结构

|                                    pktHead                                     |                    pktBody                     | 说明 |
| :----------------------------------------------------------------------------: | :--------------------------------------------: | :--- |
| `pktLen` \| `crc32` \| `sequence` \|                  `msgID`                  |                    `msgBody`                   | 客户端→服务器（网关做 router 逻辑，自行判断消息发往哪个前端服务器；或没有网关） |
| `pktLen` \| `crc32` \| `sequence` \|              `msgDirection`              | `msgID` \| `msgBody`                          | 客户端→服务器（网关不做 router 逻辑，根据 msgDirection 判断发往哪个前端服务器） |
| `pktLen` \| `crc32` \| `sequence` \|                  `msgID`                  |                    `msgBody`                   | 服务器→客户端（网关做 router 逻辑，自行判断消息发往哪个前端服务器；或没有网关） |
|                               `pktLen`                                         | `rpcType` \| `contextID` \| `msgType` \| `msgID` \| `msgBody` | 服务器内部消息 |
|                               `pktLen`                                         | `ToServiceType` \| `ToServiceID` \| `rpcType` \| `contextID` \| `msgType` \| `msgID` \| `msgBody` | 服务器内部消息 |

### WebSocket 包结构

相对于 TCP 包结构，**没有 pktLen** 字段。

---

## 服务器架构

```
          manager                               server manager (backend)
     center     matching

   login    lobby  lobby     game  game         frontend server

      gate    gate    gate                       gate server

        client  client                           client
```

---

## 待完成

- [x] 增加一种新的计时器实现
- [ ] 集成 crontab 定时调度

---

## 许可证

<!-- TODO: 补充许可证信息 -->

[返回顶部](#xtnet--go-游戏服务器网络库)

