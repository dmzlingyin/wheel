# 🛞 Wheel - 重复造轮子项目

> "虽然重复造轮子不可取，但这是我理解底层原理的最好方式"

## 项目简介

这是一个学习性质的 Go 语言项目，旨在通过手动实现各种网络协议和系统组件来深入理解底层原理。项目名称 "Wheel" 寓意着"重复造轮子"，通过亲手构建来加深对计算机科学基础概念的理解。

## 🎯 项目目标

- **学习导向**：不是为了生产使用，而是为了学习底层原理
- **从零实现**：不使用高级库，手动实现基础协议
- **渐进式学习**：从简单到复杂，逐步深入

## 📁 项目结构

```
wheel/
├── network/          # 网络协议实现
│   ├── tcp/         # TCP 客户端/服务器
│   ├── udp/         # UDP 客户端/服务器
│   ├── http/        # HTTP 服务器
│   ├── smtp/        # SMTP 邮件服务器
│   └── dns/         # DNS 客户端实现
├── go.mod
└── README.md
```

## 🚀 已实现功能

### 网络协议
- **TCP**: 基础的 TCP 服务器和客户端，展示面向连接的通信
- **UDP**: 基础的 UDP 服务器和客户端，展示无连接的快速通信
- **DNS**: DNS 客户端实现，完整实现 DNS 查询报文构建和解析过程
- **HTTP**: 自定义 HTTP 服务器，支持 GET/POST 方法
- **SMTP**: 简单邮件传输协议服务器，支持 HELO/EHLO、MAIL FROM、RCPT TO、DATA、QUIT 命令

## 📋 计划实现

### 网络协议
- [x] SMTP - 简单邮件传输协议
- [x] UDP - 用户数据报协议
- [x] DNS - 域名系统客户端
- [x] DNS - 域名系统服务器
- [ ] FTP - 文件传输协议
- [ ] DHCP - 动态主机配置协议
- [ ] WebSocket - 实时通信协议

### 操作系统相关
- [ ] 文件系统
- [ ] 进程调度
- [ ] 内存管理
- [ ] 网络栈

## 🛠️ 快速开始

### 运行 TCP 服务器
```bash
go run network/tcp/server.go
```

### 运行 TCP 客户端
```bash
go run network/tcp/client.go
```

### 运行 HTTP 服务器
```bash
go run network/http/server.go
```

### 测试 HTTP 服务器
```bash
# GET 请求
curl http://localhost:3000/

# POST 请求
curl -X POST -H "Content-Type: application/json" -d '{"option":"ping"}' http://localhost:3000/
```

### 运行 SMTP 服务器
```bash
go run network/smtp/server.go
```

### 测试 SMTP 服务器
使用 telnet 或 nc 连接测试：
```bash
telnet localhost 3000
```

SMTP 测试流程：
```
HELO example.com
MAIL FROM: <sender@example.com>
RCPT TO: <user1@example.com>
DATA
Subject: Test Email

This is a test email content.
.
QUIT
```

### 运行 UDP 服务器
```bash
go run network/udp/server.go
```

### 测试 UDP 客户端
```bash
go run network/udp/client.go
```

### 运行 DNS 客户端
```bash
go run network/dns/client.go
```
该客户端会向 DNS 服务器查询 `www.baidu.com` 的 IP 地址

## 🎓 学习价值

通过这个项目，你将学习到：

- **网络协议原理**：理解 TCP/IP 协议栈的工作方式
- **HTTP 协议细节**：掌握 HTTP 请求/响应的完整格式
- **SMTP 协议流程**：理解邮件传输协议的命令交互过程
- **UDP 无连接通信**：对比 TCP 连接和 UDP 无连接的区别
- **DNS 协议解析**：理解域名解析的完整过程和报文格式
- **并发编程**：使用 goroutine 处理并发连接
- **系统编程**：深入理解操作系统与网络的关系
- **调试技巧**：学会分析和调试网络通信问题

## 💡 设计理念

1. **简单优先**：每个实现都保持最简形式，突出核心原理
2. **注释详细**：代码中包含详细的中文注释，便于理解
3. **逐步完善**：先实现基础功能，再逐步添加特性
4. **测试驱动**：每个功能都提供可运行的测试用例

## 🤔 为什么重复造轮子？

> "我不造轮子，我只是轮子的搬运工" - 这是很多开发者的现状

但通过亲手造轮子，我们可以：
- 真正理解轮子为什么这样设计
- 知道轮子的局限性和边界
- 在遇到问题时能快速定位和修复
- 培养系统性的思维方式

## 📚 学习资源

这个项目的实现参考了：
- RFC 文档
- 《TCP/IP 详解》
- 《HTTP 权威指南》
- 各种开源项目的实现思路

## 🎉 贡献

欢迎提出建议和想法！虽然这是一个个人学习项目，但如果你有：
- 更好的实现思路
- 发现代码中的问题
- 想要学习某个特定的协议

请随时提出 issue 或讨论。

---

> **记住**：这里的每一个轮子都是为了学习而造，不是为了替代现有的优秀轮子。理解原理，方能运用自如！